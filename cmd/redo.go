package cmd

import (
	"fmt"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// redoCmd represents the redo command
var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Reapply the last migration.",
	Long:  `Reapply the last migration.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		db, dialect, err := getConnection()
		if err != nil {
			return err
		}

		source := migrate.FileMigrationSource{
			Dir: viper.GetString("database.dir"),
		}

		migrations, _, err := migrate.PlanMigration(db, dialect, source, migrate.Down, 1)
		if len(migrations) == 0 {
			fmt.Println("Nothing to do!")
			return nil
		}

		if dryRun {
			printMigration(migrations[0], migrate.Down)
			printMigration(migrations[0], migrate.Up)
		} else {
			_, err := migrate.ExecMax(db, dialect, source, migrate.Down, 1)
			if err != nil {
				return fmt.Errorf("Migration (down) failed: %s", err)
			}

			_, err = migrate.ExecMax(db, dialect, source, migrate.Up, 1)
			if err != nil {
				return fmt.Errorf("Migration (up) failed: %s", err)
			}

			fmt.Printf("Reapplied migration %s.", migrations[0].Id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(redoCmd)

	redoCmd.Flags().BoolVarP(&dryRun, "dryrun", "", false, "Don't apply migrations, just print them.")
}
