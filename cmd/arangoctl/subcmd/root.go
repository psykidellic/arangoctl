package subcmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arangoctl",
	Short: "arangoctl runs arango workflows",
}

func Execute(version string) {
	rootCmd.Version = version
	rootCmd.Execute()
}
