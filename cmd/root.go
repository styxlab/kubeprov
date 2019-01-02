package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "kubeprov",
		Short: "Kubernetes Cluster on Hetzner Cloud.",
		Long:  `Command-line interface for creating Kubernetes Clusters on Hetzner Cloud.`,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log level")

	//rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}
