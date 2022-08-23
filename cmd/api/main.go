package main

import (
	"log"

	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/application/services"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/repositories"
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

	furyHandler := server.NewHTTPServer(app, newHandlers())
	furyHandler.SetupRouter()

	return furyHandler.Run()
}

func newHandlers() server.ItemHandler {
	conn, err := repositories.GetConnectionDB()
	if err != nil {
		panic("error connecting to DB: " + err.Error())
	}

	itemRepository, err := repositories.NewItemRepository(conn)
	if err != nil {
		panic("error creating item repository: " + err.Error())
	}

	itemService := services.NewItemService(itemRepository)
	itemHandler := server.NewItemHandler(itemService)

	return itemHandler
}
