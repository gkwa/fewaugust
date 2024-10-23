package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gkwa/fewaugust/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fewaugust",
	Long:  `All software has versions. This is fewaugust's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
