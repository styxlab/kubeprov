package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
const cfgFileDefault = ".kubeprov-config"

var rootCmd = &cobra.Command{
	Use:   "kubeprov",
	Short: "CLI for provisioning a Kubernetes Cluster on Hetzner Cloud.",
	Long:  "Command-line interface for creating a Kubernetes Clusters on Hetzner Cloud.",
}

func Execute() {

	fmt.Println("execute")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	fmt.Println("init1")

	cobra.OnInitialize(initConfig)

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/" +  cfgFileDefault + ")")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log level")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	fmt.Println("init")

	if cfgFile != "" {
		fmt.Println("config flag")
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("set dir")
		setConfigDirectory()
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
    	os.Exit(1)
    }else{
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setConfigDirectory() {

	// Find the home directory
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Home-Dir:", home)

	viper.AddConfigPath(home)
	viper.SetConfigName(cfgFileDefault)
}

/*
func init() {
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
  //rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
  rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
  //rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
  rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
  viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
  viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
  viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
  viper.SetDefault("license", "MIT")
}*/
