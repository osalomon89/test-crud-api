package main

import (
	"log"

	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/repositories/mysql"
	server "github.com/osalomon89/test-crud-api/internal/infrastructure/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app, err := fury.NewWebApplication()
	if err != nil {
		return err
	}

	conn, err := mysql.GetConnectionDB()
	if err != nil {
		panic("error connecting to DB: " + err.Error())
	}
	defer conn.Close()

	serverReady := make(chan bool)
	httpServer, err := server.NewHTTPServer(app, conn, serverReady)
	if err != nil {
		panic("error creating server: " + err.Error())
	}

	httpServer.SetupRouter()
	httpServer.Run()

	return nil
}
