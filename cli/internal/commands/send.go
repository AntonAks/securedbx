package commands

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/securedbx/securedbx/cli/internal/auth"
	"github.com/securedbx/securedbx/cli/internal/config"
	"github.com/securedbx/securedbx/cli/internal/crypto"
)

var (
	sendText     string
	sendPINValue string
	sendTTL      int
)

func NewSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [file]",
		Short: "Encrypt and upload a file or text secret",
		RunE:  runSend,
	}
	cmd.Flags().StringVar(&sendText, "text", "", "Send text instead of file")
	cmd.Flags().StringVar(&sendPINValue, "pin-value", "", "PIN (4 chars a-z A-Z 0-9); auto-generated if omitted")
	cmd.Flags().IntVar(&sendTTL, "ttl", 1, "Expiry in hours (1–24)")
	return cmd
}

func runSend(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if sendText != "" && len(args) > 0 {
		return fmt.Errorf("cannot use --text with file arguments")
	}
	if len(args) > 1 {
		return fmt.Errorf("only one file at a time is supported")
	}
	if sendTTL < 1 || sendTTL > 24 {
		return fmt.Errorf("--ttl must be between 1 and 24")
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
	} else {
		return fmt.Errorf("provide a file or use --text")
	}

	pin := sendPINValue
	if pin == "" {
		pin = generatePIN()
	}

	req := api.PINUploadRequest{
		ContentType: "file",
		FileSize:    int64(len(plaintext)),
		PIN:         pin,
		FileName:    filename,
		TTL:         sendTTL * 60,
		AccessMode:  "one_time",
	}

	resp, err := client.PINUpload(ctx, req)
	if err != nil {
		return fmt.Errorf("upload init: %w", err)
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

	pr := newProgressReader(bytes.NewReader(ciphertext), int64(len(ciphertext)), "Uploading")
	if err = client.UploadToS3(ctx, resp.UploadURL, pr, int64(len(ciphertext))); err != nil {
		return fmt.Errorf("S3 upload: %w", err)
	}
	clearProgress()

	return printJSON(map[string]any{
		"file_id":    resp.FileID,
		"pin":        pin,
		"expires_at": resp.ExpiresAt,
	})
}

func generatePIN() string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const max = byte(256 - 256%len(charset))
	pin := make([]byte, 4)
	i := 0
	for i < 4 {
		var b [1]byte
		if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
			panic(fmt.Sprintf("crypto/rand failed: %v", err))
		}
		if b[0] < max {
			pin[i] = charset[int(b[0])%len(charset)]
			i++
		}
	}
	return string(pin)
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

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
