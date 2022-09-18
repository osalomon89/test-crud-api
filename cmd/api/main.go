package main

import (
	"fmt"
	"log"

	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/core/services"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/handlers"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/repositories/mysql"
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
		return fmt.Errorf("error connecting to DB: %w", err)
	}

	itemRepository, err := mysql.NewItemRepository(conn)
	if err != nil {
		return fmt.Errorf("error creating item repository: %w", err)
	}

	itemService, err := services.NewItemService(itemRepository)
	if err != nil {
		return fmt.Errorf("error creating item service: %w", err)
	}

	httpServer, err := handlers.NewHTTPServer(app, itemService)
	if err != nil {
		return fmt.Errorf("error creating server: %w", err)
	}

	httpServer.SetupRouter()

	return httpServer.Run()
}
