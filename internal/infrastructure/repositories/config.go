package repositories

import (
	"fmt"
	"os"

	//Autoload the env
	_ "github.com/joho/godotenv/autoload"
	"github.com/mercadolibre/fury_go-toolkit-config/pkg/config"
	"github.com/mercadolibre/go-meli-toolkit/gomelipass"
)

const productionEnv string = "production"

var (
	environment   = map[string]string{}
	dbHost        = "DB_HOST"
	dbUser        = "DB_USER"
	dbPass        = "DB_PASS"
	dbName        = "DB_NAME"
	dbMaxIdleConn = "DB_MAX_IDLE_CONN"
	dbMaxOpenConn = "DB_MAX_OPEN_CONN"
	dbMaxLifetime = "DB_CONN_MAX_LIFETIME"
	defaultString = ""
)

func load(env string) error {
	if env == productionEnv {
		if cfg, err := config.Load(); err != nil {
			return err
		} else {
			environment = map[string]string{
				dbHost: gomelipass.GetEnv("DB_MYSQL_DESAENV07_TESTMARKET_TESTMARKET_ENDPOINT"),
				dbUser: cfg.GetString(dbUser, defaultString),
				dbPass: gomelipass.GetEnv("DB_MYSQL_DESAENV07_TESTMARKET_TESTMARKET_WPROD"),
				dbName: cfg.GetString(dbName, defaultString),
			}
			return nil
		}
	}

	loadENVConfigs()
	return nil
}

func loadENVConfigs() {
	scope := os.Getenv("SCOPE")

	if scope == "" {
		environment = map[string]string{
			dbHost:        os.Getenv(dbHost),
			dbUser:        os.Getenv(dbUser),
			dbPass:        os.Getenv(dbPass),
			dbName:        os.Getenv(dbName),
			dbMaxIdleConn: os.Getenv(dbMaxIdleConn),
			dbMaxOpenConn: os.Getenv(dbMaxOpenConn),
			dbMaxLifetime: os.Getenv(dbMaxLifetime),
		}
	}
}

// dbConnectionURL returns db connecction url
func dbConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", environment[dbUser], environment[dbPass], environment[dbHost], environment[dbName])
}
