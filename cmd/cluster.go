package cmd

import (
	"github.com/spf13/cobra"
)

var (
	clusterCmd = &cobra.Command {
		Use:     "cluster",
		Short:   "Manages the Cluster",
	}
)

func init() {
	clusterCmd.AddCommand(createCmd)
}