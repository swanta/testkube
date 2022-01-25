package tests

import (
	"os"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"github.com/kubeshop/testkube/pkg/ui"
	"github.com/spf13/cobra"
)

func NewTestExecutionsCmd() *cobra.Command {
	var (
		limit int
		tags  []string
	)

	cmd := &cobra.Command{
		Use:     "executions [testName]",
		Aliases: []string{"el"},
		Short:   "Gets tests executions list",
		Long:    `Gets tests executions list, can be filtered by test name`,
		Run: func(cmd *cobra.Command, args []string) {
			ui.Logo()

			var testName string
			if len(args) > 0 {
				testName = args[0]
			}

			client, _ := common.GetClient(cmd)

			executions, err := client.ListTestExecutions(testName, limit, tags)
			ui.ExitOnError("getting tests executions list", err)

			ui.Table(executions, os.Stdout)
			ui.NL()

		},
	}

	cmd.Flags().IntVar(&limit, "limit", 1000, "max number of records to return")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "comma separated list of tags: --tags tag1,tag2,tag3")

	return cmd
}
