package cmd

import (
	"os"

	"github.com/bradcypert/deckard/lib/db"
	"github.com/bradcypert/deckard/lib/migrations"
	"github.com/spf13/cobra"
)

func upFunc(args []string) {
	bindVarsFromConfig()

	database := db.Database{
		Dbname:   cmdDatabaseName,
		Port:     cmdDatabasePort,
		Password: cmdDatabasePassword,
		User:     cmdDatabaseUser,
		Host:     cmdDatabaseHost,
		Driver:   cmdDatabaseDriver,
		IsSilent: cmdIsSilent,
	}

	if len(args) < 1 {
		// get all migrations in current folder.
		migrations := migrations.Migrations{cmdIsSilent}
		migration := migrations.FindInPath(cmdInputDir, true)
		database.RunUp(migration, cmdSteps)
	} else {
		// TODO: What if we have more args?
	}
}

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Runs one or more \"up\" migrations.",
	Long:  `Runs one or more \"up\" migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		configOverwriter(cmd)
		upFunc(args)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	addDatabaseFlags(upCmd)

	upCmd.Flags().IntVarP(&cmdSteps,
		"steps",
		"s",
		-1,
		"The number of up migrations you'd like to run.")

	dir, _ := os.Getwd()
	upCmd.Flags().StringVarP(&cmdInputDir,
		"inputDir",
		"i",
		dir,
		"Directory which contains the migrations")
}
