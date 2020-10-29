package cmd

import (
	"fmt"
	"os"

	"github.com/isayme/go-wstunnel/wstunnel"
	"github.com/spf13/cobra"
)

var versionFlag bool

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "wstunnel version")
}

var rootCmd = &cobra.Command{
	Use: "wstunnel",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			wstunnel.PrintVersion()
			os.Exit(0)
		}
	},
}

// Execute run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
