package main

import (
	"log"

	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/core/services"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/repositories/mysql"
	server "github.com/osalomon89/test-crud-api/internal/infrastructure/server"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server/handler"
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

func newHandlers() handler.ItemHandler {
	conn, err := mysql.GetConnectionDB()
	if err != nil {
		panic("error connecting to DB: " + err.Error())
	}

	itemRepository, err := mysql.NewItemRepository(conn)
	if err != nil {
		panic("error creating item repository: " + err.Error())
	}

	itemService, err := services.NewItemService(itemRepository)
	if err != nil {
		panic("error creating item service: " + err.Error())
	}

	itemHandler, err := handler.NewItemHandler(itemService)
	if err != nil {
		panic("error creating item handler: " + err.Error())
	}

	return itemHandler
}
