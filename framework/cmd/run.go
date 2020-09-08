package cmd

import (
	"fmt"
	"strings"

	"github.com/julianGoh17/simple-e2e/framework/operations"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/spf13/cobra"
)

var (
	stages string
	test   string
)

// NewRunCmd returns the run command as a cobra object to be interacted with
func NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run through a set of or all the stages in a test",
		Long:  `Run all the steps in a specified test or just a specific set of stages from that test.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			util.ConfigureGlobalLogLevel(verbosity)
			// TODO: add documentation that says that each run should be stateless
			controller, err := operations.NewController()
			if err != nil {
				return err
			}
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
	runCmd.Flags().StringVarP(&test, "test", "t", "", "The name of the test to run. Do not need to pass in file extension.")
	runCmd.Flags().StringVarP(&stages, "stages", "s", "", `A comma separated list of stages to run from that test.
For example to only run 'stage1' from a test, add '-s stage1' to your command.
	`)
	runCmd.MarkFlagRequired("test")
	rootCmd.AddCommand(runCmd)
}
