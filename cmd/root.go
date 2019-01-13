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
	//rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log level")
	//rootCmd.PersistentFlags().BoolP("port", "p", false, "port for logging (default: 9090)")
	//rootCmd.PersistentFlags().BoolP("logs", "l", false, "show logs (default: run commands)")

	//rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loggingCmd)
}

func Execute() {

	err := rootCmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}

}
