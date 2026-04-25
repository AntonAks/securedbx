package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/securedbx/securedbx/cli/internal/commands"
)

func main() {
	root := &cobra.Command{
		Use:               "sdbx",
		Short:             "Zero-knowledge file sharing CLI for securedbx.com",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	root.AddCommand(commands.NewSendCmd())
	root.AddCommand(commands.NewReceiveCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
