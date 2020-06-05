package example

import (
	"github.com/spf13/cobra"
)

// Cmd cmd example
var Cmd = &cobra.Command{
	Use:   "example:test",
	Short: "示例",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}
