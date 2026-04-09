package commands

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/securedbx/securedbx/cli/internal/auth"
	"github.com/securedbx/securedbx/cli/internal/config"
	"github.com/securedbx/securedbx/cli/internal/crypto"
)

var (
	sendText      string
	sendPIN       bool
	sendPINValue  string
	sendPassword  bool
	sendPassValue string
	sendMulti     bool
	sendTTL       string
	sendJSON      bool
)

func NewSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [file ...]",
		Short: "Encrypt and upload a file or text secret",
		RunE:  runSend,
	}
	cmd.Flags().StringVar(&sendText, "text", "", "Send text instead of file")
	cmd.Flags().BoolVar(&sendPIN, "pin", false, "Use PIN mode")
	cmd.Flags().StringVar(&sendPINValue, "pin-value", "", "PIN explicitly (4-8 digits)")
	cmd.Flags().BoolVar(&sendPassword, "password", false, "Password-protect the file")
	cmd.Flags().StringVar(&sendPassValue, "password-value", "", "Password explicitly")
	cmd.Flags().BoolVar(&sendMulti, "multi", false, "Multi-access mode (use with --password)")
	cmd.Flags().StringVar(&sendTTL, "ttl", "1h", "TTL: 5m, 1h, 12h, 24h, 3d, 7d")
	cmd.Flags().BoolVar(&sendJSON, "json", false, "Output result as JSON")
	return cmd
}

func runSend(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if sendText != "" && len(args) > 0 {
		return fmt.Errorf("cannot use --text with file arguments")
	}
	if sendPIN && sendPassword {
		return fmt.Errorf("cannot use --pin and --password together")
	}

	unauthClient := api.New(config.BaseURL(), config.CFSecret(), "")
	apiKey, err := auth.EnsureValidKey(ctx, unauthClient)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	client := api.New(config.BaseURL(), config.CFSecret(), apiKey)

	var plaintext []byte
	var filename string

	if sendText != "" {
		plaintext = []byte(sendText)
		filename = "secret.txt"
	} else if len(args) == 1 {
		plaintext, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		filename = filepath.Base(args[0])
	} else if len(args) > 1 {
		var entries []crypto.FileEntry
		for _, path := range args {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read %s: %w", path, err)
			}
			entries = append(entries, crypto.FileEntry{Name: filepath.Base(path), Data: data})
		}
		plaintext, err = crypto.BundleFiles(entries)
		if err != nil {
			return fmt.Errorf("bundle files: %w", err)
		}
		filename = "archive.zip"
	} else {
		return fmt.Errorf("provide at least one file or use --text")
	}

	if sendPIN {
		return sendPINMode(ctx, client, plaintext, filename)
	}
	return sendURLMode(ctx, client, plaintext, filename)
}

func sendURLMode(ctx context.Context, client *api.Client, plaintext []byte, filename string) error {
	accessMode := "one_time"
	if sendMulti {
		accessMode = "multi"
	}

	var ciphertext, key []byte
	var encryptedKey, saltB64 string
	var err error

	if sendPassword {
		password, err := readPassword("Enter password: ")
		if err != nil {
			return err
		}

		salt, err := crypto.GenerateSalt()
		if err != nil {
			return fmt.Errorf("generate salt: %w", err)
		}
		saltB64 = base64.StdEncoding.EncodeToString(salt)

		wrapKey := crypto.DeriveKeyFromPassword(password, salt)

		dataKey := make([]byte, 32)
		if _, err = io.ReadFull(rand.Reader, dataKey); err != nil {
			return fmt.Errorf("generate data key: %w", err)
		}

		ciphertext, err = crypto.EncryptWithKey(plaintext, dataKey)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}

		wrappedKey, err := crypto.EncryptKey(dataKey, wrapKey)
		if err != nil {
			return fmt.Errorf("wrap key: %w", err)
		}
		encryptedKey = base64.StdEncoding.EncodeToString(wrappedKey)
		accessMode = "multi"
	} else {
		ciphertext, key, err = crypto.Encrypt(plaintext)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	}

	req := api.UploadInitRequest{
		ContentType:  "file",
		FileSize:     int64(len(ciphertext)),
		TTL:          sendTTL,
		AccessMode:   accessMode,
		Salt:         saltB64,
		EncryptedKey: encryptedKey,
	}
	resp, err := client.UploadInit(ctx, req)
	if err != nil {
		return fmt.Errorf("upload init: %w", err)
	}

	if err = client.UploadToS3(ctx, resp.UploadURL, ciphertext); err != nil {
		return fmt.Errorf("upload to S3: %w", err)
	}

	if sendJSON {
		if sendPassword {
			fmt.Printf(`{"mode":"password","file_id":%q,"expires_at":%d}`+"\n", resp.FileID, resp.ExpiresAt)
		} else {
			shareURL := buildShareURL(resp.FileID, key, filename)
			fmt.Printf(`{"mode":"url","url":%q,"file_id":%q,"expires_at":%d}`+"\n", shareURL, resp.FileID, resp.ExpiresAt)
		}
		return nil
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	expiry := time.Unix(resp.ExpiresAt, 0).Format("2006-01-02 15:04 MST")

	if sendPassword {
		fmt.Printf("\n%s Password-protected file ready!\n", green.Render("✓"))
		fmt.Printf("  File ID: %s  |  Expires: %s\n", resp.FileID, expiry)
		fmt.Printf("  %s\n\n", yellow.Render("Share the File ID and password separately."))
	} else {
		shareURL := buildShareURL(resp.FileID, key, filename)
		fmt.Printf("\n%s Done! Share this link (one-time, %s):\n", green.Render("✓"), sendTTL)
		fmt.Printf("  %s\n\n", shareURL)
		fmt.Printf("  %s\n\n", yellow.Render("⚠  Link contains the key. Keep it private."))
	}
	return nil
}

func sendPINMode(ctx context.Context, client *api.Client, plaintext []byte, filename string) error {
	pin := sendPINValue
	if pin == "" {
		pin = generatePIN()
	}

	req := api.PINUploadRequest{
		ContentType: "file",
		FileSize:    int64(len(plaintext)),
		PIN:         pin,
		FileName:    filename,
		TTL:         sendTTL,
		AccessMode:  "one_time",
	}

	resp, err := client.PINUpload(ctx, req)
	if err != nil {
		return fmt.Errorf("PIN upload init: %w", err)
	}

	saltBytes, err := hexDecodeString(resp.Salt)
	if err != nil {
		return fmt.Errorf("decode salt: %w", err)
	}
	key := crypto.DeriveKeyFromPIN(pin, saltBytes)

	ciphertext, err := crypto.EncryptWithKey(plaintext, key)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	if err = client.UploadToS3(ctx, resp.UploadURL, ciphertext); err != nil {
		return fmt.Errorf("S3 upload: %w", err)
	}

	expiry := time.Unix(resp.ExpiresAt, 0).Format("2006-01-02 15:04 MST")

	if sendJSON {
		fmt.Printf(`{"mode":"pin","file_id":%q,"pin":%q,"expires_at":%d}`+"\n", resp.FileID, pin, resp.ExpiresAt)
		return nil
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	fmt.Printf("\n%s Done!\n", green.Render("✓"))
	fmt.Printf("  File ID: %s  |  PIN: %s  |  Expires: %s\n\n", resp.FileID, pin, expiry)
	return nil
}

func buildShareURL(fileID string, key []byte, filename string) string {
	keyStr := crypto.KeyToBase64URL(key)
	return fmt.Sprintf("https://securedbx.com/#/download?id=%s&key=%s&name=%s", fileID, keyStr, urlEncode(filename))
}

func readPassword(prompt string) (string, error) {
	if sendPassValue != "" {
		return sendPassValue, nil
	}
	fmt.Fprint(os.Stderr, prompt)
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	return string(b), err
}

func generatePIN() string {
	b := make([]byte, 1)
	rand.Reader.Read(b)
	return fmt.Sprintf("%04d", int(b[0])%10000)
}

func urlEncode(s string) string {
	return strings.NewReplacer(" ", "%20", "#", "%23").Replace(s)
}

func hexDecodeString(s string) ([]byte, error) {
	b := make([]byte, len(s)/2)
	for i := 0; i < len(s)-1; i += 2 {
		n, err := fmt.Sscanf(s[i:i+2], "%02x", &b[i/2])
		if err != nil || n != 1 {
			return nil, fmt.Errorf("invalid hex at position %d", i)
		}
	}
	return b, nil
}
