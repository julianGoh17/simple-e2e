package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	stages string
	test   string
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a specified test or stages of that test",
		Long:  `Run all the steps in a specified test or just a specific set of stages from that specified test.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			stage := []string{}
			if stages != "" {
				stage = strings.Split(stages, ",")
			}
			return controller.RunTest(test+".yaml", stage...)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	// TODO: Better description
	runCmd.Flags().StringVarP(&test, "test", "t", "", "The name of the test to run. Do not need to pass in file extension. ")
	runCmd.Flags().StringVarP(&stages, "stages", "s", "", "A comma separated list of stages to run from that test.")
	runCmd.MarkFlagRequired("test")
}
