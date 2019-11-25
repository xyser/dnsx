package api

import (
	"dnsx/api"
	"dnsx/model/dao"
	"github.com/spf13/cobra"
)

// ServerCmd http server
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "start http api server",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		dao.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}
