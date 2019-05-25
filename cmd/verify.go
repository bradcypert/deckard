// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"deckard/db"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

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
		bindVarsFromConfig()
		var migration db.Migration
		queries := make([]db.Query, 0)

		if strings.HasSuffix(args[0],".up.sql") {
			contents, _ := ioutil.ReadFile(args[0])
			queries = append(queries, db.Query{
				Name:  "",
				Value: string(contents),
			})
		}

		migration = db.Migration {
			Queries: queries,
		}

		postgres := db.Postgres{
			Dbname: cmdDatabaseName,
			Port: cmdDatabasePort,
			Password: cmdDatabasePassword,
			User: cmdDatabaseUser,
			Host: cmdDatabaseHost,
		}

		postgres.Verify(migration)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&cmdDatabaseConfigSelector,
		"key",
		"k",
		"",
		"The database key to use from the YAML config provided in the configFile argument.")

	verifyCmd.Flags().StringVarP(&cmdDatabaseHost,
		"host",
		"t",
		"",
		"The host for the database you'd like to apply the up migrations to.")

	verifyCmd.Flags().StringVarP(&cmdDatabaseName,
		"database",
		"d",
		"",
		"The database name that you'd like to apply the up migrations to")

	verifyCmd.Flags().StringVarP(&cmdDatabaseUser,
		"user",
		"u",
		"",
		"The user you'd like to connect to the database as.")

	verifyCmd.Flags().StringVarP(&cmdDatabasePassword,
		"password",
		"a",
		"",
		"The password for the database user that you're applying migrations as.")

	verifyCmd.Flags().IntVarP(&cmdDatabasePort,
		"port",
		"p",
		0,
		"The port that the database you're targeting runs on.")
}
