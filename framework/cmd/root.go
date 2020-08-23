package cmd

import (
	"fmt"
	"os"

	"github.com/julianGoh17/simple-e2e/framework/operations"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	home    = "/home/e2e"
	rootCmd = &cobra.Command{
		Use:   "Simple-E2E",
		Short: "A modular and configurable testing infrastructure",
		Long: `Simple-E2E is a testing library aimed at making more modular and easier. 
This application allows users to break down tests into stages and steps to 
run a set of stages or an entire test. Furthermore, Simple-E2E provides a
framework to easily create new tests from exisiting steps.`,
	}
	controller *operations.Controller
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "Julian Goh", "author name for copyright attribution")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("license", "apache")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
