package mysql

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB //nolint:gochecknoglobals

func GetConnectionDB() (*sqlx.DB, error) {
	var err error
	env := os.Getenv("GO_ENVIRONMENT")

	if db == nil {
		if err := load(env); err != nil {
			return nil, fmt.Errorf("### CONFIGS ERROR: %w", err)
		}

		db, err = sqlx.Connect("mysql", dbConnectionURL())
		if err != nil {
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
			return nil, fmt.Errorf("### DB ERROR: %w", err)
		}
	}

	if env != productionEnv {
		if err := autoMigrate(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func autoMigrate(db *sqlx.DB) error {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
		return fmt.Errorf("error instantiating migration: %w", err)
	}

	dbMigration, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migration",
		"mysql",
		driver,
	)

	if err != nil {
		fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
		return fmt.Errorf("error instantiating migration: %w", err)
	}

	if err := dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error executing migration: %w", err)
	}

	return nil
}
