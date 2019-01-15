package cmd

import (
	"database/sql"
	"fmt"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"gopkg.in/gorp.v1"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var dialects = map[string]gorp.Dialect{
	"sqlite3":  gorp.SqliteDialect{},
	"postgres": gorp.PostgresDialect{},
	"mysql":    gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"},
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
	}

	viper.SetDefault("database.dir", "migrations")
	viper.SetDefault("database.table", "migrations")

	//add also parser from env variables for easier docker access
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if viper.GetString("database.schema") != "" {
		migrate.SetSchema(viper.GetString("database.schema"))
	}

	migrate.SetTable(viper.GetString("database.table"))
}

func getConnection() (*sql.DB, string, error) {

	// Make sure we only accept dialects that were compiled in.
	dialect := viper.GetString("database.dialect")
	_, exists := dialects[dialect]
	if !exists {
		return nil, "", fmt.Errorf("Unsupported dialect: %s", dialect)
	}

	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.protocol"),
		viper.GetString("database.address"),
		viper.GetString("database.dbname"),
	)
	viper.Set("database.datasource", dsn)

	db, err := sql.Open(dialect, viper.GetString("database.datasource"))
	if err != nil {
		return nil, "", fmt.Errorf("Cannot connect to database: %s", err)
	}

	return db, dialect, nil
}
