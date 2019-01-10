package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var name string

var templateContent = `
-- +migrate Up

-- +migrate Down
`
var tpl *template.Template

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new a database migration.",
	Long:  `Create a new a database migration.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(viper.GetString("database.dir")); os.IsNotExist(err) {
			return err
		}

		fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
		pathName := path.Join(viper.GetString("database.dir"), fileName)
		f, err := os.Create(pathName)

		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		if err := tpl.Execute(f, nil); err != nil {
			return err
		}

		fmt.Printf("Created migration %s", pathName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the migration.")
	_ = newCmd.MarkFlagRequired("name")

	tpl = template.Must(template.New("new_migration").Parse(templateContent))
}
