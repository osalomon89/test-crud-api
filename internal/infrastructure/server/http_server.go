package server

import (
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
)

type HTTPServer interface {
	SetupRouter()
	Run() error
}

type httpServer struct {
	ItemHandler ItemHandler
	App         *fury.Application
}

func NewHTTPServer(app *fury.Application, handler ItemHandler) HTTPServer {
	return &httpServer{
		ItemHandler: handler,
		App:         app,
	}
}

func (handler *httpServer) SetupRouter() {
	api := handler.App.Router.Group("/v1/items")
	{
		api.Get("/{id}", handler.ItemHandler.GetItemByID)
		api.Post("/", handler.ItemHandler.CreateItem)
	}
}

func (handler *httpServer) Run() error {
	return handler.App.Run()
}
