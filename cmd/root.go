package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgFileDefault = ".kubeprov/config"

var rootCmd = &cobra.Command{
	Use:   "kubeprov",
	Short: "CLI for provisioning a Kubernetes Cluster on Hetzner Cloud.",
	Long:  "Command-line interface for creating a Kubernetes Clusters on Hetzner Cloud.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log level")
}

// initConfig reads in config file if present
func initConfig() {

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(cfgFileDefault)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
    	//os.Exit(1)
    }else{
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}