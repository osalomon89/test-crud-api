package config

import (
	"fmt"
	"os"

	//Autoload the env
	_ "github.com/joho/godotenv/autoload"
)

var (
	environment   = map[string]string{}
	port          = "PORT"
	dbHost        = "DB_HOST"
	dbPort        = "DB_PORT"
	dbUser        = "DB_USER"
	dbPass        = "DB_PASS"
	dbName        = "DB_NAME"
	dbMaxIdleConn = "DB_MAX_IDLE_CONN"
	dbMaxOpenConn = "DB_MAX_OPEN_CONN"
	dbMaxLifetime = "DB_CONN_MAX_LIFETIME"
)

func Load() {
	environment = map[string]string{
		port:          os.Getenv(port),
		dbHost:        os.Getenv(dbHost),
		dbPort:        os.Getenv(dbPort),
		dbUser:        os.Getenv(dbUser),
		dbPass:        os.Getenv(dbPass),
		dbName:        os.Getenv(dbName),
		dbMaxIdleConn: os.Getenv(dbMaxIdleConn),
		dbMaxOpenConn: os.Getenv(dbMaxOpenConn),
		dbMaxLifetime: os.Getenv(dbMaxLifetime),
	}
}

// DBConnectionURL returns db connecction url
func DBConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", environment[dbUser], environment[dbPass], environment[dbHost], environment[dbPort], environment[dbName])
}

// Port return server port
func Port() string {
	return environment[port]
}
