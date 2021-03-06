package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cmdDatabaseConfigSelector string
var cmdDatabasePassword string
var cmdDatabaseHost string
var cmdDatabasePort int
var cmdDatabaseUser string
var cmdDatabaseName string
var cmdDatabaseDriver string
var cmdInputDir string
var cmdIsSilent bool
var cmdSteps int
var cmdSSLConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deckard",
	Short: "Deckard helps manage database migrations.",
	Long: `Deckard helps manage database migrations. Deckard has spawned from the need to decouple migrations
from any particular application. Deck'. For example:

Deckard helps you manage and create migrations for multiple databases!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "A YAML config file containing the database definition per the Deckard specification (default is $HOME/.deckard.yaml). See example.deckard.yml in the github repo for an example.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".deckard" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".deckard")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
