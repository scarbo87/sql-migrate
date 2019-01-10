package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type statusRow struct {
	Id        string
	Migrated  bool
	AppliedAt time.Time
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status.",
	Long:  `Show migration status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, dialect, err := getConnection()
		if err != nil {
			return err
		}

		source := migrate.FileMigrationSource{
			Dir: viper.GetString("database.dir"),
		}
		migrations, err := source.FindMigrations()
		if err != nil {
			return err
		}

		records, err := migrate.GetMigrationRecords(db, dialect)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Migration", "Applied"})
		table.SetColWidth(60)

		rows := make(map[string]*statusRow)

		for _, m := range migrations {
			rows[m.Id] = &statusRow{
				Id:       m.Id,
				Migrated: false,
			}
		}

		for _, r := range records {
			if rows[r.Id] == nil {
				//ui.Warn(fmt.Sprintf("Could not find migration file: %v", r.Id))
				continue
			}

			rows[r.Id].Migrated = true
			rows[r.Id].AppliedAt = r.AppliedAt
		}

		for _, m := range migrations {
			if rows[m.Id] != nil && rows[m.Id].Migrated {
				table.Append([]string{
					m.Id,
					rows[m.Id].AppliedAt.String(),
				})
			} else {
				table.Append([]string{
					m.Id,
					"no",
				})
			}
		}

		table.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
