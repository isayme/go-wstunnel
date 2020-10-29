package cmd

import (
	"github.com/isayme/go-wstunnel/cmd/local"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(localCmd)
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Run local of wstunnel",
	Run: func(cmd *cobra.Command, args []string) {
		local.Run()
	},
}
