package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// caneCmd represents the cane command
var caneCmd = &cobra.Command{
	Use:   "cane",
	Short: "Ponder mysteries of the Horadrim",
	Long: `Ponder mysteries of the Horadrim`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stay a while and listen...")
	},
}

func init() {
	rootCmd.AddCommand(caneCmd)
}
