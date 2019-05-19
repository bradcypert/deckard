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
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var upCmdDatabaseConfigSelector string
var upCmdDatabasePassword string
var upCmdDatabaseHost string
var upCmdDatabasePort int
var upCmdDatabaseUser string
var upCmdDatabaseName string


// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Runs one or more \"up\" migrations.",
	Long: `Runs one or more \"up\" migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		var migration db.Migration
		queries := make([]db.Query, 3)

		dir, _ := os.Getwd()

		if len(args) < 1 {
			// get all migrations in current folder.
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range files {
				if strings.HasSuffix(file.Name(),".up.sql") {
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
				Dbname: upCmdDatabaseName,
				Port: upCmdDatabasePort,
				Password: upCmdDatabasePassword,
				User: upCmdDatabaseUser,
				Host:upCmdDatabaseHost,
			}

			postgres.RunUp(migration)
		} else {

		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().StringVarP(&upCmdDatabaseConfigSelector,
		"dbKey",
		"k",
		"",
		"The database key to use from the YAML config provided in the configFile argument.")

	upCmd.Flags().StringVarP(&upCmdDatabaseHost,
		"host",
		"t",
		"",
		"The host for the database you'd like to apply the up migrations to.")

	upCmd.Flags().StringVarP(&upCmdDatabaseName,
		"database",
		"d",
		"",
		"The database name that you'd like to apply the up migrations to")

	upCmd.Flags().StringVarP(&upCmdDatabaseUser,
		"user",
		"u",
		"",
		"The user you'd like to connect to the database as.")

	upCmd.Flags().StringVarP(&upCmdDatabasePassword,
		"password",
		"a",
		"",
		"The password for the database user that you're applying migrations as.")

	upCmd.Flags().IntVarP(&upCmdDatabasePort,
		"port",
		"p",
		0,
		"The port that the database you're targeting runs on.")


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
