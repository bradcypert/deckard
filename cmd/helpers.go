package cmd

import "github.com/spf13/viper"

func bindVarsFromConfig() {
	if cmdDatabaseConfigSelector != "" {
		cmdDatabasePort = viper.GetInt(cmdDatabaseConfigSelector+".port")
		cmdDatabasePassword = viper.GetString(cmdDatabaseConfigSelector+".password")
		cmdDatabaseUser = viper.GetString(cmdDatabaseConfigSelector+".user")
		cmdDatabaseHost = viper.GetString(cmdDatabaseConfigSelector+".host")
		cmdDatabaseName = viper.GetString(cmdDatabaseConfigSelector+".database")
	}
}
