package cmd

import (
	"fmt"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// skipCmd represents the skip command
var skipCmd = &cobra.Command{
	Use:   "skip",
	Short: "Set the database level to the most recent version available, without actually running the migrations.",
	Long:  `Set the database level to the most recent version available, without actually running the migrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		db, dialect, err := getConnection()
		if err != nil {
			return err
		}

		source := migrate.FileMigrationSource{
			Dir: viper.GetString("database.dir"),
		}

		n, err := migrate.SkipMax(db, dialect, source, migrate.Up, limit)
		if err != nil {
			return fmt.Errorf("Migration failed: %s", err)
		}

		if n == 1 {
			fmt.Println("Skipped 1 migration")
		} else {
			fmt.Printf("Skipped %d migrations\n", n)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(skipCmd)

	skipCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Max number of migrations to apply.")
}
