package repositories

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB //nolint:gochecknoglobals

func GetConnectionDB() (*gorm.DB, error) {
	var err error
	env := os.Getenv("GO_ENVIRONMENT")

	if db == nil {
		if err := load(env); err != nil {
			return nil, fmt.Errorf("### CONFIGS ERROR: %w", err)
		}

		db, err = gorm.Open(mysql.Open(dbConnectionURL()), &gorm.Config{})
		if err != nil {
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
			return nil, fmt.Errorf("### DB ERROR: %w", err)
		}
	}

	if env != productionEnv {
		if err := migrate(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(Item{}, Photo{})
	if err != nil {
		fmt.Printf("########## MIGRATE ERROR: " + err.Error() + " #############")

		return fmt.Errorf("### MIGRATE ERROR: %w", err)
	}

	return nil
}
