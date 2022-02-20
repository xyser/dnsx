package cmd

import (
	"fmt"
	"os"

	"github.com/xyser/dnsx/cmd/api"
	"github.com/xyser/dnsx/cmd/secret"
	"github.com/xyser/dnsx/pkg/config"
	"github.com/xyser/dnsx/pkg/log"

	"github.com/spf13/cobra"
)

// RootCmd RootCmd
var RootCmd = &cobra.Command{
	Use:              "dnsx",
	Short:            "DNS Server",
	Long:             "dns service of dnsx.",
	TraverseChildren: true,
}

func init() {
	cobra.OnInitialize(onInitialize)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&config.CfgFile, "c", "config/dnsx.yaml", "config file (default is config.yaml)")

	RootCmd.AddCommand(api.ServerCmd)
	RootCmd.AddCommand(secret.Cmd)

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("\033[1;30;42m[error]\033[0m commands execute error: %s", err)
		os.Exit(1)
	}
}

func onInitialize() {
	// 初始化依赖包
	config.Init()
	log.Init()
}
