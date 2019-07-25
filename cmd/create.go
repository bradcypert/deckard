package cmd

import (
	"github.com/bradcypert/deckard/lib/migrations"
	"github.com/spf13/cobra"
)

var createOutputDir string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new migration",
	Long: `Creates a new migration to be ran by the database. Does not actually run
any migration with this command, however.

Use: deckard create add_login_date_to_users`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Create(createOutputDir, args[0])
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
