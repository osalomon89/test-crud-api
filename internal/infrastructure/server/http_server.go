package server

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/osalomon89/test-crud-api/internal/core/services"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/repositories/mysql"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server/handler"
)

type HTTPServer interface {
	SetupRouter()
	Run()
}

type httpServer struct {
	ItemHandler handler.ItemHandler
	App         *fury.Application
	ServerReady chan bool
}

func NewHTTPServer(app *fury.Application, conn *sqlx.DB, serverReady chan bool) (HTTPServer, error) {
	handler, err := newItemHandler(conn)
	if err != nil {
		return nil, err
	}

	return &httpServer{
		ItemHandler: handler,
		App:         app,
		ServerReady: serverReady,
	}, nil
}

func (h *httpServer) SetupRouter() {
	api := h.App.Router.Group("/v1/items")
	{
		api.Get("/{id}", h.ItemHandler.GetItemByID)
		api.Post("/", h.ItemHandler.CreateItem)
	}
}

func (h *httpServer) Run() {
	go func() {
		if err := h.App.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	if h.ServerReady != nil {
		h.ServerReady <- true
	}
}

func newItemHandler(conn *sqlx.DB) (handler.ItemHandler, error) {
	itemRepository, err := mysql.NewItemRepository(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating item repository: %w", err)
	}

	itemService, err := services.NewItemService(itemRepository)
	if err != nil {
		return nil, fmt.Errorf("error creating item service: %w", err)
	}

	itemHandler, err := handler.NewItemHandler(itemService)
	if err != nil {
		return nil, fmt.Errorf("error creating item handler: %w", err)
	}

	return itemHandler, nil
}
