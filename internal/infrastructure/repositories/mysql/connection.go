package mysql

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
		if err := migrate(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func migrate(db *sqlx.DB) error {
	var itemsSchema = `
	CREATE TABLE IF NOT EXISTS items (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		code varchar(191) DEFAULT NULL,
		title longtext,
		description longtext,
		price bigint(20) DEFAULT NULL,
		stock bigint(20) DEFAULT NULL,
		item_type longtext,
		leader tinyint(1) DEFAULT NULL,
		leader_level longtext,
		status longtext,
		created_at datetime(3) DEFAULT NULL,
		updated_at datetime(3) DEFAULT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY code (code)
	  );`

	_, err := db.Exec(itemsSchema)
	if err != nil {
		fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
		return fmt.Errorf("### MIGRATION ERROR: %w", err)
	}

	var photosSchema = `CREATE TABLE IF NOT EXISTS photos (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		path longtext,
		item_id bigint(20) unsigned DEFAULT NULL,
		created_at datetime(3) DEFAULT NULL,
		updated_at datetime(3) DEFAULT NULL,
		PRIMARY KEY (id),
		KEY fk_items_photos (item_id),
		CONSTRAINT fk_items_photos FOREIGN KEY (item_id) REFERENCES items (id)
	  );`

	_, err = db.Exec(photosSchema)
	if err != nil {
		fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
		return fmt.Errorf("### MIGRATION ERROR: %w", err)
	}

	return nil
}
