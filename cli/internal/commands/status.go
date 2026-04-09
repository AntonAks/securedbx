package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/securedbx/securedbx/cli/internal/api"
	"github.com/securedbx/securedbx/cli/internal/auth"
	"github.com/securedbx/securedbx/cli/internal/config"
)

var statusJSON bool

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <file-id>",
		Short: "Show file metadata",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatus,
	}
	cmd.Flags().BoolVar(&statusJSON, "json", false, "Output as JSON")
	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	fileID := args[0]

	unauthClient := api.New(config.BaseURL(), config.CFSecret(), "")
	apiKey, err := auth.EnsureValidKey(ctx, unauthClient)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	client := api.New(config.BaseURL(), config.CFSecret(), apiKey)

	meta, err := client.Metadata(ctx, fileID)
	if err != nil {
		return err
	}

	if statusJSON {
		b, _ := json.MarshalIndent(meta, "", "  ")
		fmt.Println(string(b))
		return nil
	}

	expiry := time.Unix(meta.ExpiresAt, 0).Format("2006-01-02 15:04 MST")
	avail := "yes"
	if !meta.Available {
		avail = "no (already downloaded)"
	}

	bold := lipgloss.NewStyle().Bold(true)
	fmt.Printf("\n%s %s\n", bold.Render("File ID:"), meta.FileID)
	fmt.Printf("%s %s\n", bold.Render("Type:"), meta.ContentType)
	fmt.Printf("%s %d bytes\n", bold.Render("Size:"), meta.FileSize)
	fmt.Printf("%s %s\n", bold.Render("Expires:"), expiry)
	fmt.Printf("%s %s\n", bold.Render("Available:"), avail)
	fmt.Printf("%s %s\n\n", bold.Render("Access mode:"), meta.AccessMode)
	return nil
}
