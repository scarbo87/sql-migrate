# sql-migrate-cobra

> SQL Schema migration tool for [Go](http://golang.org/). Based on [sql-migrate](https://github.com/rubenv/sql-migrate).

## Features

* See [sql-migrate](https://github.com/rubenv/sql-migrate).

## Installation

To install the library and command line program, use the following:

```bash
go get -v github.com/scarbo87/sql-migrate-cobra/...
```

## Usage

### As a standalone tool

```
$ sql-migrate-cobra --help
Usage:
  sql-migrate-cobra [command]

Available Commands:
  down        Undo a database migration.
  help        Help about any command
  new         Create a new a database migration.
  redo        Reapply the last migration.
  skip        Set the database level to the most recent version available, without actually running the migrations.
  status      Show migration status.
  up          Migrates the database to the most recent version available.

Flags:
      --config string   config file (default is ./config.toml)
  -h, --help            help for sql-migrate-cobra

```

Each command requires a configuration file (which defaults to `config.toml`, but can be specified with the `--config` flag). Example:

```toml
[database]
dialect = "mysql"
username = "username"
password = "password"
protocol = "tcp"
address = "127.0.0.1:3306"
dbname = "database"
dir = "migrations"
table = "migrations"
```

Run with env override:

```bash
DATABASE_USERNAME=root DATABASE_PASSWORD=root DATABASE_DBNAME=dbname sql-migrate-cobra`
```

The `table` setting is optional and will default to `migrations`.

Use the `--help` flag in combination with any of the commands to get an overview of its usage:

```
$ sql-migrate-cobra up --help
Migrates the database to the most recent version available.

Usage:
  sql-migrate-cobra up [flags]

Flags:
      --dry-run     Don't apply migrations, just print them.
  -h, --help        help for up
      --limit int   Max number of migrations to apply.

Global Flags:
      --config string   config file (default is ./config.toml)
```

The `new` command creates a new empty migration template using the following pattern `<current time>-<name>.sql`.

The `up` command applies all available migrations. By contrast, `down` will only apply one migration by default. This behavior can be changed for both by using the `-limit` parameter.

The `redo` command will unapply the last migration and reapply it. This is useful during development, when you're writing migrations.

Use the `status` command to see the state of the applied migrations:

```bash
$ sql-migrate-cobra status
+---------------+-----------------------------------------+
|   MIGRATION   |                 APPLIED                 |
+---------------+-----------------------------------------+
| 1_initial.sql | 2014-09-13 08:19:06.788354925 +0000 UTC |
| 2_record.sql  | no                                      |
+---------------+-----------------------------------------+
```

### MySQL Caveat

If you are using MySQL, you must append `?parseTime=true` to the `datasource` configuration.


## Writing migrations
Migrations are defined in SQL files, which contain a set of SQL statements. Special comments are used to distinguish up and down migrations.

```sql
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE people (id int);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
```

You can put multiple statements in each block, as long as you end them with a semicolon (`;`).

You can alternatively set up a separator string that matches an entire line by setting `sqlparse.LineSeparator`. This
can be used to imitate, for example, MS SQL Query Analyzer functionality where commands can be separated by a line with
contents of `GO`. If `sqlparse.LineSeparator` is matched, it will not be included in the resulting migration scripts.

If you have complex statements which contain semicolons, use `StatementBegin` and `StatementEnd` to indicate boundaries:

```sql
-- +migrate Up
CREATE TABLE people (id int);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION do_something()
returns void AS $$
DECLARE
  create_query text;
BEGIN
  -- Do something here
END;
$$
language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
DROP FUNCTION do_something();
DROP TABLE people;
```

The order in which migrations are applied is defined through the filename: sql-migrate will sort migrations based on their name. It's recommended to use an increasing version number or a timestamp as the first part of the filename.

Normally each migration is run within a transaction in order to guarantee that it is fully atomic. However some SQL commands (for example creating an index concurrently in PostgreSQL) cannot be executed inside a transaction. In order to execute such a command in a migration, the migration can be run using the `notransaction` option:

```sql
-- +migrate Up notransaction
CREATE UNIQUE INDEX people_unique_id_idx CONCURRENTLY ON people (id);

-- +migrate Down
DROP INDEX people_unique_id_idx;
```

## License

This library is distributed under the [MIT](LICENSE) license.
