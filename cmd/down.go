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
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Runs one or more \"down\" migrations.",
	Long: `Runs one or more \"down\" migrations.
These migrations are likely destructive. Please use caution when executing deckard down.

Deckard can be instructed to run all down migrations or specific ones.

Running All:
Example:
deckard down

Running One:
Example:
deckard down 1558294955321
# or
deckard down add_users_to_other_users
`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmdDatabaseConfigSelector != "" {
			cmdDatabasePort = viper.GetInt(cmdDatabaseConfigSelector+".port")
			cmdDatabasePassword = viper.GetString(cmdDatabaseConfigSelector+".password")
			cmdDatabaseUser = viper.GetString(cmdDatabaseConfigSelector+".user")
			cmdDatabaseHost = viper.GetString(cmdDatabaseConfigSelector+".host")
			cmdDatabaseName = viper.GetString(cmdDatabaseConfigSelector+".database")
		}
		var migration db.Migration
		queries := make([]db.Query, 0)

		if len(args) < 1 {
			// get all migrations in current folder.
			files, err := ioutil.ReadDir(cmdInputDir)
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range files {
				if strings.HasSuffix(file.Name(),".down.sql") {
					contents, _ := ioutil.ReadFile(file.Name())
					queries = append(queries, db.Query{
						Name:  file.Name(),
						Value: string(contents),
					})
				}
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

			postgres.RunDown(migration)
		} else {

		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().StringVarP(&cmdDatabaseConfigSelector,
		"dbKey",
		"k",
		"",
		"The database key to use from the YAML config provided in the configFile argument.")

	downCmd.Flags().StringVarP(&cmdDatabaseHost,
		"host",
		"t",
		"",
		"The host for the database you'd like to apply the down migrations to.")

	downCmd.Flags().StringVarP(&cmdDatabaseName,
		"database",
		"d",
		"",
		"The database name that you'd like to apply the down migrations to")

	downCmd.Flags().StringVarP(&cmdDatabaseUser,
		"user",
		"u",
		"",
		"The user you'd like to connect to the database as.")

	downCmd.Flags().StringVarP(&cmdDatabasePassword,
		"password",
		"a",
		"",
		"The password for the database user that you're applying migrations as.")

	downCmd.Flags().IntVarP(&cmdDatabasePort,
		"port",
		"p",
		0,
		"The port that the database you're targeting runs on.")

	dir, _ := os.Getwd()
	downCmd.Flags().StringVarP(&cmdInputDir,
		"inputDir",
		"i",
		dir,
		"Directory which contains the migrations")
}
