package cmd

import (
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Undo a database migration.",
	Long:  `Undo a database migration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return applyMigrations(migrate.Down, dryRun, limit)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
	downCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't apply migrations, just print them.")
}
