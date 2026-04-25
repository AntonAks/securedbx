package commands

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/securedbx/securedbx/cli/internal/auth"
	"github.com/securedbx/securedbx/cli/internal/config"
	"github.com/securedbx/securedbx/cli/internal/crypto"
)

var (
	receivePINValue string
	receiveOutput   string
)

func NewReceiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive <file-id>",
		Short: "Download and decrypt a file or text secret",
		Args:  cobra.ExactArgs(1),
		RunE:  runReceive,
	}
	cmd.Flags().StringVar(&receivePINValue, "pin-value", "", "PIN (prompts if omitted)")
	cmd.Flags().StringVar(&receiveOutput, "output", ".", "Output directory or file path")
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

	return receiveWithPIN(ctx, client, args[0])
}

func receiveWithPIN(ctx context.Context, client *api.Client, fileID string) error {
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
		ciphertext, err := client.DownloadFromS3(ctx, verifyResp.DownloadURL, func(read, total int64) {
			printProgress("Downloading", read, total)
		})
		clearProgress()
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
		return printJSON(map[string]any{
			"type":    "text",
			"content": string(plaintext),
		})
	}

	if filename == "" {
		filename = "download"
	}
	savedPath, err := saveOutput(plaintext, filename)
	if err != nil {
		return err
	}
	return printJSON(map[string]any{
		"type": "file",
		"path": savedPath,
	})
}

func saveOutput(plaintext []byte, filename string) (string, error) {
	outputPath := receiveOutput

	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		outputPath = filepath.Join(outputPath, filename)
	}

	if strings.HasSuffix(filename, ".zip") {
		entries, err := crypto.UnbundleFiles(plaintext)
		if err == nil && len(entries) > 0 {
			outDir := filepath.Clean(filepath.Dir(outputPath))
			for _, e := range entries {
				destPath := filepath.Join(outDir, filepath.Base(e.Name))
				if !strings.HasPrefix(destPath, outDir+string(os.PathSeparator)) {
					return "", fmt.Errorf("invalid zip entry name: %q", e.Name)
				}
				if err = os.WriteFile(destPath, e.Data, 0o600); err != nil {
					return "", fmt.Errorf("write %s: %w", e.Name, err)
				}
			}
			return outDir, nil
		}
	}

	if err = os.WriteFile(outputPath, plaintext, 0o600); err != nil {
		return "", fmt.Errorf("write output: %w", err)
	}
	return outputPath, nil
}
