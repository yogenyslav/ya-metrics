package database

import (
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// RunMigrations applies last migrations to database.
func RunMigration(db DB, dialect string) error {
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	conn, err := db.SqlDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := goose.Up(conn, "migrations"); err != nil {
		return err
	}

	return nil
}
