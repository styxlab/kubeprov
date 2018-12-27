package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kubeprov",
	Short: "CLI for provisioning a Kubernetes Cluster on Hetzner Cloud",
	Long:  "Command-line interface for creating a Kubernetes Clusters on Hetzner Cloud.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	/*if debug, err := rootCmd.PersistentFlags().GetBool("debug"); err != nil && debug {
		pkg.RenderProgressBars = false
	} else {
		pkg.RenderProgressBars = true
	}*/

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
