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
		Short: "Print the version number",
		Long:  `Print the version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Simple-E2E binary version: v0.1")
		},
	}
}
