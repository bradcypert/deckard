package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var createOutputDir string

func createFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()

	return err
}

func makeTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new migration",
	Long: `Creates a new migration to be ran by the database. Does not actually run
any migration with this command, however.

Use: deckard create add_login_date_to_users`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		// Add in the / suffix if it wasn't added by the user
		if len(createOutputDir) > 0 && !strings.HasSuffix(createOutputDir, "/") {
			createOutputDir += "/"
		}

		filepath := createOutputDir + makeTimestamp() + "__" + args[0]
		upError := createFile(filepath + ".up.sql")
		downError := createFile(filepath + ".down.sql")

		fmt.Printf("Created file %s\n", filepath + ".up.sql")
		fmt.Printf("Created file %s\n", filepath + ".down.sql")

		if upError != nil {
			log.Fatal(upError)
		}

		if downError != nil {
			log.Fatal(downError)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&createOutputDir, "outputDir", "o", "", "Output directory to write the migration to, defaults to current directory.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
