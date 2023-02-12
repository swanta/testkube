package commands

import (
	"github.com/spf13/cobra"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common/validator"
	"github.com/kubeshop/testkube/pkg/ui"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Shows version and build info",
		Long:    `Shows version and build info`,
		Run: func(cmd *cobra.Command, args []string) {
			client, _ := common.GetClient(cmd)
			info, err := client.GetServerInfo()
			if err != nil {
				info.Version = info.Version + " " + err.Error()
			}

			ui.Logo()
			ui.Info("Client Version", common.Version)
			ui.Info("Server Version", info.Version)
			ui.Info("Commit", common.Commit)
			ui.Info("Built by", common.BuiltBy)
			ui.Info("Build date", common.Date)

		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validator.PersistentPreRunVersionCheck(cmd, common.Version)
		},
	}
}
