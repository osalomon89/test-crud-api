package handlers

import (
	"fmt"

	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/core/ports"
)

type HTTPServer interface {
	SetupRouter()
	Run() error
}

type httpServer struct {
	itemHandler *itemHandler
	app         *fury.Application
}

func NewHTTPServer(app *fury.Application, itemService ports.ItemService) (HTTPServer, error) {
	itemHandler, err := newItemHandler(itemService)
	if err != nil {
		return nil, fmt.Errorf("error creating item handler")
	}

	return &httpServer{
		itemHandler: itemHandler,
		app:         app,
	}, nil
}

func (server *httpServer) SetupRouter() {
	api := server.app.Router.Group("/v1/items")
	{
		api.Get("/{id}", server.itemHandler.GetItemByID)
		api.Post("/", server.itemHandler.CreateItem)
	}
}

func (handler *httpServer) Run() error {
	return handler.app.Run()
}
