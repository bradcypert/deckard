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
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var downCmdDatabaseConfigSelector string
var downCmdDatabasePassword string
var downCmdDatabaseHost string
var downCmdDatabasePort int
var downCmdDatabaseUser string
var downCmdDatabaseName string
var downCmdInputDir string

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
		var migration db.Migration
		queries := make([]db.Query, 3)

		if len(args) < 1 {
			// get all migrations in current folder.
			files, err := ioutil.ReadDir(downCmdInputDir)
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
				Dbname: downCmdDatabaseName,
				Port: downCmdDatabasePort,
				Password: downCmdDatabasePassword,
				User: downCmdDatabaseUser,
				Host: downCmdDatabaseHost,
			}

			postgres.RunDown(migration)
		} else {

		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().StringVarP(&downCmdDatabaseConfigSelector,
		"dbKey",
		"k",
		"",
		"The database key to use from the YAML config provided in the configFile argument.")

	downCmd.Flags().StringVarP(&downCmdDatabaseHost,
		"host",
		"t",
		"",
		"The host for the database you'd like to apply the down migrations to.")

	downCmd.Flags().StringVarP(&downCmdDatabaseName,
		"database",
		"d",
		"",
		"The database name that you'd like to apply the down migrations to")

	downCmd.Flags().StringVarP(&downCmdDatabaseUser,
		"user",
		"u",
		"",
		"The user you'd like to connect to the database as.")

	downCmd.Flags().StringVarP(&downCmdDatabasePassword,
		"password",
		"a",
		"",
		"The password for the database user that you're applying migrations as.")

	downCmd.Flags().IntVarP(&downCmdDatabasePort,
		"port",
		"p",
		0,
		"The port that the database you're targeting runs on.")

	dir, _ := os.Getwd()
	downCmd.Flags().StringVarP(&downCmdInputDir,
		"inputDir",
		"i",
		dir,
		"Directory which contains the migrations")


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
