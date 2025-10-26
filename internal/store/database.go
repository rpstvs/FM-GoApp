package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {

	godotenv.Load(".env")

	dbString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", dbString)

	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	return db, nil

}

func MigrateFs(db *sql.DB, migrationFs fs.FS, dir string) error {
	goose.SetBaseFS(migrationFs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)

	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
