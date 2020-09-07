package cmd

import (
	"os"

	"github.com/julianGoh17/simple-e2e/framework/operations"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	allFlag bool
)

// NewListCmd returns the list command as a cobra object to be interacted with
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists the container names and IDs",
		Long:  `Lists the container names and IDs running on the host's daemon`,
		RunE: func(cmd *cobra.Command, args []string) error {
			util.ConfigureGlobalLogLevel(verbosity)
			// TODO: Add ability to controll to show all containers/running containers and show state of containers
			controller, err := operations.NewController()
			if err != nil {
				return err
			}
			namesAndIDs, err := controller.GetContainerNamesAndIDs()
			if err != nil {
				return err
			}

			table := getTable(namesAndIDs)
			table.Render()

			return nil
		},
	}
}

func initListCmd(rootCmd, listCmd *cobra.Command) {
	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Determine whether to list all containers. Default value 'false'.")

	rootCmd.AddCommand(listCmd)
}

func getTable(namesAndIDs map[string]string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Container Name", "Container ID"})
	table.SetBorder(false)

	for name, id := range namesAndIDs {
		table.Append([]string{name, id})
	}

	return table

}
