package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func bindVarsFromConfig() {
	if cmdDatabaseConfigSelector != "" {
		cmdDatabasePort = viper.GetInt(cmdDatabaseConfigSelector+".port")
		cmdDatabasePassword = viper.GetString(cmdDatabaseConfigSelector+".password")
		cmdDatabaseUser = viper.GetString(cmdDatabaseConfigSelector+".user")
		cmdDatabaseHost = viper.GetString(cmdDatabaseConfigSelector+".host")
		cmdDatabaseName = viper.GetString(cmdDatabaseConfigSelector+".database")
		cmdDatabaseDriver = viper.GetString(cmdDatabaseConfigSelector+".driver")
	}
}

func addDatabaseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cmdDatabaseConfigSelector,
		"key",
		"k",
		"",
		"The database key to use from the YAML config provided in the configFile argument.")

	cmd.Flags().StringVarP(&cmdDatabaseHost,
		"host",
		"t",
		"",
		"The host for the database you'd like to apply the migrations to.")

	cmd.Flags().StringVarP(&cmdDatabaseName,
		"database",
		"d",
		"",
		"The database name that you'd like to apply the migrations to")

	cmd.Flags().StringVarP(&cmdDatabaseUser,
		"user",
		"u",
		"",
		"The user you'd like to connect to the database as.")

	cmd.Flags().StringVarP(&cmdDatabasePassword,
		"password",
		"a",
		"",
		"The password for the database user that you're applying migrations as.")

	cmd.Flags().IntVarP(&cmdDatabasePort,
		"port",
		"p",
		0,
		"The port that the database you're targeting runs on.")

	cmd.Flags().StringVarP(&cmdDatabaseDriver,
		"driver",
		"r",
		"",
		"The database driver for connecting to the database. Valid options are: [mysql, postgres]")
}