package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Current version and track of kubeprov.
const version = "v0.0.1"
const track = "DEV"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "prints the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubeprov " + version + "-" + track)
	},
}
