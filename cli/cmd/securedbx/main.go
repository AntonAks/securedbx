package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/securedbx/securedbx/cli/internal/commands"
)

func main() {
	root := &cobra.Command{
		Use:   "securedbx",
		Short: "Zero-knowledge file sharing CLI for securedbx.com",
	}

	root.AddCommand(commands.NewSendCmd())
	root.AddCommand(commands.NewReceiveCmd())
	root.AddCommand(commands.NewStatusCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
