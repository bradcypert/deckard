package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/bradcypert/deckard/lib/db"
	"github.com/bradcypert/deckard/lib/migrations"
	"github.com/spf13/cobra"
)

func verifyFunc(args []string) {
	bindVarsFromConfig()
	var migration migrations.Migration
	queries := make([]migrations.Query, 0)

	if strings.HasSuffix(args[0], ".up.sql") {
		contents, _ := ioutil.ReadFile(args[0])
		queries = append(queries, migrations.Query{
			Name:  "",
			Value: string(contents),
		})
	}

	migration = migrations.Migration{
		Queries: queries,
	}

	database := db.Database{
		Dbname:   cmdDatabaseName,
		Port:     cmdDatabasePort,
		Password: cmdDatabasePassword,
		User:     cmdDatabaseUser,
		Host:     cmdDatabaseHost,
		Driver:   cmdDatabaseDriver,
		IsSilent: cmdIsSilent,
	}

	database.Verify(migration)
}

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that a given migration exists in the database.",
	Long: `Verifies that a given migration exists in the database. Due to the nature of Deckard, only UP migrations
are stored in the database, so verifying a down migration is going to always fail.

Provide the path to the migration to verify.

Example:
deckard verify ./migrations/1234_add_login_date_to_users.up.sql`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verifyFunc(args)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	addDatabaseFlags(verifyCmd)
}
