package api

import (
	"github.com/dingdayu/dnsx/api"
	"github.com/dingdayu/dnsx/model/mysql"
	"github.com/dingdayu/dnsx/pkg/config"
	"github.com/dingdayu/dnsx/pkg/validate"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		mysql.Init()
		// load mode
		if config.GetString("app.mode") == "release" {
			gin.SetMode(gin.ReleaseMode)
		}
		// customize validator
		binding.Validator = validate.GinValidator()
	},
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}

func init() {
	ServerCmd.PersistentFlags().BoolVar(&authoritative, "authoritative", true, "type: 是否启用权威服务; default: true")
	_ = viper.BindPFlag("dns.authoritative", ServerCmd.PersistentFlags().Lookup("author"))
}
