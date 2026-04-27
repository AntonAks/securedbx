package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/securedbx/securedbx/cli/internal/commands"
	"github.com/securedbx/securedbx/cli/internal/updater"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:               "sdbx",
		Short:             "Zero-knowledge file sharing CLI for securedbx.com",
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	root.AddCommand(commands.NewSendCmd())
	root.AddCommand(commands.NewReceiveCmd())
	root.AddCommand(commands.NewVersionCmd(version))

	// Check cache before running (no latency); refreshes cache in background if stale.
	latest := updater.Check(version)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if latest != "" {
		fmt.Fprintf(os.Stderr, "\nUpdate available: %s → %s\n", version, latest)
		fmt.Fprintf(os.Stderr, "Download: https://github.com/AntonAks/sdbx/releases/latest\n")
	}
}
