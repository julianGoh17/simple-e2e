package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initVersionCmd(rootCmd, versionCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}

// NewVersionCmd returns the version command as a cobra command object
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of the Simple-E2E",
		Long:  `Print the current version of the installed Simple-E2E binary`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Simple-E2E binary version: v0.1")
		},
	}
}
