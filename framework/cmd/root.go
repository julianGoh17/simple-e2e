package cmd

import (
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbosity string
	config    = util.GlobalConfig{}
)

// NewRootCmd returns the root cli as an object to be interacted with
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "Simple-E2E",
		Short: "A modular and configurable testing infrastructure",
		Long: `Simple-E2E is a testing library aimed at making more modular and easier. 
		This application allows users to break down tests into stages and steps to 
		run a set of stages or an entire test. Furthermore, Simple-E2E provides a
		framework to easily create new tests from exisiting steps.`,
	}
}

// InitRootCmd configures the root comand with the basic information and adds the other commands into the binary
func InitRootCmd(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "", `Increase the verbosity of the binary by passing in one of the following levels:
	info: Will log basic events (default)
	debug: Will increase logging level to show what step and stage is being called
	trace: Will increase logging level to show debug level + more
		`)

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("license", "apache")

	versionCmd := NewVersionCmd()
	initVersionCmd(rootCmd, versionCmd)

	runCmd := NewRunCmd()
	initRunCmd(rootCmd, runCmd)

	listCmd := NewListCmd()
	initListCmd(rootCmd, listCmd)
}
