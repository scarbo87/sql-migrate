package cmd

import (
	"errors"
	"fmt"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	dryRun  bool
	limit   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sql-migrate-cobra",
	Short: "",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("database.dialect") == "" {
			return errors.New("No dialect specified")
		}

		if viper.GetString("database.datasource") == "" {
			return errors.New("No data source specified")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.toml)")
}

func applyMigrations(dir migrate.MigrationDirection, dryRun bool, limit int) error {

	db, dialect, err := getConnection()
	if err != nil {
		return err
	}

	source := migrate.FileMigrationSource{
		Dir: viper.GetString("database.dir"),
	}

	if dryRun {
		migrations, _, err := migrate.PlanMigration(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Cannot plan migration: %s", err)
		}

		for _, m := range migrations {
			printMigration(m, dir)
		}
	} else {
		n, err := migrate.ExecMax(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Migration failed: %s", err)
		}

		if n == 1 {
			fmt.Println("Applied 1 migration")
		} else {
			fmt.Printf("Applied %d migrations\n", n)
		}
	}

	return nil
}

func printMigration(m *migrate.PlannedMigration, dir migrate.MigrationDirection) {
	if dir == migrate.Up {
		fmt.Printf("==> Would apply migration %s (up)\n", m.Id)
		for _, q := range m.Up {
			fmt.Println(q)
		}
	} else if dir == migrate.Down {
		fmt.Printf("==> Would apply migration %s (down)\n", m.Id)
		for _, q := range m.Down {
			fmt.Println(q)
		}
	} else {
		panic("Not reached")
	}
}
