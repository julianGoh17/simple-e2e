package cmd

import (
	"fmt"
	"strings"

	"github.com/julianGoh17/simple-e2e/framework/operations"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/spf13/cobra"
)

var (
	stages    string
	test      string
	verbosity string
	config    = util.GlobalConfig{}
)

// NewRunCmd returns the run command as a cobra object to be interacted with
func NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run a specified test or stages of that test",
		Long:  `Run all the steps in a specified test or just a specific set of stages from that specified test.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			util.ConfigureGlobalLogLevel(verbosity)
			controller = operations.NewController()
			stage := []string{}
			if stages != "" {
				stage = strings.Split(stages, ",")
			}
			testPath := fmt.Sprintf("%s/%s.yaml", config.GetOrDefault(util.TestDirEnv), test)
			return controller.RunTest(testPath, stage...)
		},
	}
}

func initRunCmd(rootCmd, runCmd *cobra.Command) {
	// TODO: Better description
	runCmd.Flags().StringVarP(&test, "test", "t", "", "The name of the test to run. Do not need to pass in file extension.")
	runCmd.Flags().StringVarP(&stages, "stages", "s", "", "A comma separated list of stages to run from that test.")
	runCmd.Flags().StringVarP(&verbosity, "verbosity", "v", "", `Increase the verbosity of the binary by passing in one of the following levels:
		info: Will log basic events (default)
		debug: Will increase logging level to show what step and stage is being called
		trace: Will increase logging level to show debug level + more
	`)

	runCmd.MarkFlagRequired("test")
	rootCmd.AddCommand(runCmd)
}
