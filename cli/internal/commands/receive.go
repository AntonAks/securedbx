package commands

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/securedbx/securedbx/cli/internal/auth"
	"github.com/securedbx/securedbx/cli/internal/config"
	"github.com/securedbx/securedbx/cli/internal/crypto"
)

var (
	receivePIN       string
	receivePINValue  string
	receivePassword  bool
	receivePassValue string
	receiveOutput    string
	receiveJSON      bool
)

func NewReceiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive [url]",
		Short: "Download and decrypt a file or text secret",
		RunE:  runReceive,
	}
	cmd.Flags().StringVar(&receivePIN, "pin", "", "PIN mode: provide file ID")
	cmd.Flags().StringVar(&receivePINValue, "pin-value", "", "PIN (prompts if omitted)")
	cmd.Flags().BoolVar(&receivePassword, "password", false, "Password-protected file")
	cmd.Flags().StringVar(&receivePassValue, "password-value", "", "Password (prompts if omitted)")
	cmd.Flags().StringVar(&receiveOutput, "output", ".", "Output directory or file path")
	cmd.Flags().BoolVar(&receiveJSON, "json", false, "Output metadata as JSON")
	return cmd
}

func runReceive(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	unauthClient := api.New(config.BaseURL(), config.CFSecret(), "")
	apiKey, err := auth.EnsureValidKey(ctx, unauthClient)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	client := api.New(config.BaseURL(), config.CFSecret(), apiKey)

	if receivePIN != "" {
		return receiveWithPIN(ctx, client)
	}
	if receivePassword {
		if len(args) == 0 {
			return fmt.Errorf("provide file ID: securedbx receive --password <file-id>")
		}
		return receiveWithPassword(ctx, client, args[0])
	}
	if len(args) == 0 {
		return fmt.Errorf("provide a URL or use --pin / --password")
	}
	return receiveURLMode(ctx, client, args[0])
}

func receiveURLMode(ctx context.Context, client *api.Client, rawURL string) error {
	fileID, key, filename, err := parseShareURL(rawURL)
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}

	dlResp, err := client.Download(ctx, fileID)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	var plaintext []byte

	if dlResp.ContentType == "text" {
		ct, err := base64.StdEncoding.DecodeString(dlResp.EncryptedText)
		if err != nil {
			return fmt.Errorf("decode encrypted text: %w", err)
		}
		plaintext, err = crypto.Decrypt(ct, key)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		ciphertext, err := client.DownloadFromS3(ctx, dlResp.DownloadURL)
		if err != nil {
			return fmt.Errorf("S3 download: %w", err)
		}
		plaintext, err = crypto.Decrypt(ciphertext, key)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	}

	if dlResp.AccessMode == "one_time" {
		if err = client.Confirm(ctx, fileID); err != nil {
			fmt.Fprintf(os.Stderr, "warning: confirm failed: %v\n", err)
		}
	}

	if dlResp.ContentType == "text" {
		fmt.Println(string(plaintext))
		return nil
	}
	return saveOutput(plaintext, filename)
}

func receiveWithPIN(ctx context.Context, client *api.Client) error {
	fileID := receivePIN

	_, err := client.PINInitiate(ctx, api.PINInitiateRequest{FileID: fileID})
	if err != nil {
		return fmt.Errorf("initiate PIN session: %w", err)
	}

	pin := receivePINValue
	if pin == "" {
		fmt.Fprint(os.Stderr, "Enter PIN: ")
		b, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return err
		}
		pin = string(b)
	}

	verifyResp, err := client.PINVerify(ctx, api.PINVerifyRequest{FileID: fileID, PIN: pin})
	if err != nil {
		return fmt.Errorf("PIN verify: %w", err)
	}

	saltBytes, err := hex.DecodeString(verifyResp.Salt)
	if err != nil {
		return fmt.Errorf("decode salt: %w", err)
	}
	key := crypto.DeriveKeyFromPIN(pin, saltBytes)

	var plaintext []byte
	if verifyResp.ContentType == "text" {
		ct, err := base64.StdEncoding.DecodeString(verifyResp.EncryptedText)
		if err != nil {
			return fmt.Errorf("decode encrypted text: %w", err)
		}
		plaintext, err = crypto.Decrypt(ct, key)
		if err != nil {
			return fmt.Errorf("decrypt text: %w", err)
		}
	} else {
		ciphertext, err := client.DownloadFromS3(ctx, verifyResp.DownloadURL)
		if err != nil {
			return fmt.Errorf("S3 download: %w", err)
		}
		plaintext, err = crypto.Decrypt(ciphertext, key)
		if err != nil {
			return fmt.Errorf("decrypt file: %w", err)
		}
	}

	filename := verifyResp.FileName
	if verifyResp.ContentType == "text" || filename == "secret.txt" {
		fmt.Println(string(plaintext))
		return nil
	}
	if filename == "" {
		filename = "download"
	}
	return saveOutput(plaintext, filename)
}

func receiveWithPassword(ctx context.Context, client *api.Client, fileID string) error {
	meta, err := client.Metadata(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}
	if meta.Salt == "" || meta.EncryptedKey == "" {
		return fmt.Errorf("file is not password-protected")
	}

	password := receivePassValue
	if password == "" {
		fmt.Fprint(os.Stderr, "Enter password: ")
		b, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return err
		}
		password = string(b)
	}

	salt, err := base64.StdEncoding.DecodeString(meta.Salt)
	if err != nil {
		return fmt.Errorf("decode salt: %w", err)
	}
	wrapKey := crypto.DeriveKeyFromPassword(password, salt)

	encryptedKey, err := base64.StdEncoding.DecodeString(meta.EncryptedKey)
	if err != nil {
		return fmt.Errorf("decode encrypted key: %w", err)
	}
	dataKey, err := crypto.DecryptKey(encryptedKey, wrapKey)
	if err != nil {
		return fmt.Errorf("decrypt key (wrong password?): %w", err)
	}

	dlResp, err := client.Download(ctx, fileID)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	var plaintext []byte
	if dlResp.ContentType == "text" {
		ct, err := base64.StdEncoding.DecodeString(dlResp.EncryptedText)
		if err != nil {
			return fmt.Errorf("decode encrypted text: %w", err)
		}
		plaintext, err = crypto.DecryptWithKey(ct, dataKey)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		ciphertext, err := client.DownloadFromS3(ctx, dlResp.DownloadURL)
		if err != nil {
			return fmt.Errorf("S3 download: %w", err)
		}
		plaintext, err = crypto.DecryptWithKey(ciphertext, dataKey)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	}

	return saveOutput(plaintext, "download")
}

func parseShareURL(rawURL string) (fileID string, key []byte, filename string, err error) {
	hashIdx := strings.Index(rawURL, "#")
	if hashIdx < 0 {
		return "", nil, "", fmt.Errorf("URL missing # fragment")
	}
	fragment := rawURL[hashIdx+1:]

	qIdx := strings.Index(fragment, "?")
	if qIdx < 0 {
		return "", nil, "", fmt.Errorf("URL missing query parameters")
	}
	queryStr := fragment[qIdx+1:]

	params, err := url.ParseQuery(queryStr)
	if err != nil {
		return "", nil, "", fmt.Errorf("parse query: %w", err)
	}

	fileID = params.Get("id")
	keyStr := params.Get("key")
	filename = params.Get("name")

	if fileID == "" || keyStr == "" {
		return "", nil, "", fmt.Errorf("URL missing id or key parameters")
	}

	key, err = crypto.Base64URLToKey(keyStr)
	if err != nil {
		return "", nil, "", fmt.Errorf("decode key: %w", err)
	}
	return fileID, key, filename, nil
}

func saveOutput(plaintext []byte, filename string) error {
	outputPath := receiveOutput

	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		if filename == "" {
			filename = "download"
		}
		outputPath = filepath.Join(outputPath, filename)
	}

	if strings.HasSuffix(filename, ".zip") {
		entries, err := crypto.UnbundleFiles(plaintext)
		if err == nil && len(entries) > 0 {
			for _, e := range entries {
				destPath := filepath.Join(filepath.Dir(outputPath), e.Name)
				if err = os.WriteFile(destPath, e.Data, 0o644); err != nil {
					return fmt.Errorf("write %s: %w", e.Name, err)
				}
			}
			green := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
			fmt.Printf("\n%s Extracted %d file(s) to %s\n\n", green.Render("✓"), len(entries), filepath.Dir(outputPath))
			return nil
		}
	}

	if err = os.WriteFile(outputPath, plaintext, 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	fmt.Printf("\n%s Saved to %s\n\n", green.Render("✓"), outputPath)
	return nil
}
