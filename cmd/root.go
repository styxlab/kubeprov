package cmd

import (
	"fmt"
	"os"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubeprov",
	Short: "Kubernetes Cluster on Hetzner Cloud.",
	Long:  "Command-line interface for creating a Kubernetes Clusters on Hetzner Cloud.",
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log level")

	log.SetFlags(0)
	log.SetPrefix("kubeprov: ")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
