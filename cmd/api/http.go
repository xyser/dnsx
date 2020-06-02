package api

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dingdayu/dnsx/api"
	"github.com/dingdayu/dnsx/model/dao"
)

var authoritative bool

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

func init() {
	ServerCmd.PersistentFlags().BoolVar(&authoritative, "authoritative", true, "type: 是否启用权威服务; default: true")
	_ = viper.BindPFlag("dns.authoritative", ServerCmd.PersistentFlags().Lookup("author"))
}
