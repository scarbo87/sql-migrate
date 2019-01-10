package cmd

import (
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrates the database to the most recent version available.",
	Long:  `Migrates the database to the most recent version available.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return applyMigrations(migrate.Up, dryRun, limit)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().IntVar(&limit, "limit", 0, "Max number of migrations to apply.")
	upCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't apply migrations, just print them.")
}
