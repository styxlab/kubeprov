package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// to be populated by linker later
	version = "v0.0.1"
	track = "DEV"

	versionCmd = &cobra.Command	{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "prints the current version",
		Run: getVersion,
	}
)

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Println("kubeprov " + version + "-" + track)
}
