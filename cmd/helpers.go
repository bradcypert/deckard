package cmd

import (
	"fmt"
	"log"

	"github.com/bradcypert/deckard/lib/migrations"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func configOverwriter(cmd *cobra.Command) {
	overwriteConfigField(cmd, "port")
	overwriteConfigField(cmd, "password")
	overwriteConfigField(cmd, "user")
	overwriteConfigField(cmd, "host")
	overwriteConfigField(cmd, "database")
	overwriteConfigField(cmd, "driver")
	overwriteConfigField(cmd, "sslmode")
}

func overwriteConfigField(cmd *cobra.Command, field string) {
	viperField := fmt.Sprintf("%s.%s", cmdDatabaseConfigSelector, field)

	if err := viper.BindPFlag(viperField, cmd.Flags().Lookup(field)); err != nil {
		log.Printf("Cannot overwrite %s value from config file. %s\n", field, err.Error())
	}
}

func bindVarsFromConfig() {
	if cmdDatabaseConfigSelector != "" {
		cmdDatabasePort = viper.GetInt(cmdDatabaseConfigSelector + ".port")
		cmdDatabasePassword = viper.GetString(cmdDatabaseConfigSelector + ".password")
		cmdDatabaseUser = viper.GetString(cmdDatabaseConfigSelector + ".user")
		cmdDatabaseHost = viper.GetString(cmdDatabaseConfigSelector + ".host")
		cmdDatabaseName = viper.GetString(cmdDatabaseConfigSelector + ".database")
		cmdDatabaseDriver = viper.GetString(cmdDatabaseConfigSelector + ".driver")
	}
}

func addDatabaseFlags(cmd *cobra.Command) {
	addSilentFlag(cmd)

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

func addSilentFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&cmdIsSilent,
		"silent",
		false,
		"Supress non-fatal log messages")
}

func ReverseQuerySlice(a []migrations.Query) []migrations.Query {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}
