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
		migrations := migrations.Migrations{
			cmdIsSilent,
		}
		migrations.Create(createOutputDir, args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	AddSilentFlag(createCmd)
	createCmd.Flags().StringVarP(&createOutputDir, "outputDir", "o", "", "Output directory to write the migration to, defaults to current directory.")
}
