package main

import (
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/handlers"
)

func setupRouter(app *fury.Application, itemHandler handlers.ItemHandler) {
	api := app.Router.Group("/v1/items")
	{
		api.Get("/{id}", itemHandler.GetItemByID)
		api.Post("/", itemHandler.CreateItem)
	}
}
