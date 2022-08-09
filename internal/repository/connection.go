package repository

import (
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/config"
	"github.com/osalomon89/test-crud-api/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB //nolint:gochecknoglobals

func GetConnectionDB() (*gorm.DB, error) {
	var err error

	if db == nil {
		db, err = gorm.Open(mysql.Open(config.DBConnectionURL()), &gorm.Config{})
		if err != nil {
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")

			return nil, fmt.Errorf("### DB ERROR: %w", err)
		}
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.SetupJoinTable(&model.User{}, "Items", &model.UserItem{})
	if err != nil {
		fmt.Printf("########## JOIN ERROR: " + err.Error() + " #############")

		return fmt.Errorf("### JOIN ERROR: %w", err)
	}

	err = db.AutoMigrate(model.User{}, model.Item{}, model.Photo{})
	if err != nil {
		fmt.Printf("########## MIGRATE ERROR: " + err.Error() + " #############")

		return fmt.Errorf("### MIGRATE ERROR: %w", err)
	}

	return nil
}
