package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osalomon89/test-crud-api/internal/handler"
	"github.com/osalomon89/test-crud-api/internal/repository"
	"github.com/osalomon89/test-crud-api/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Server represents server
type Server struct {
	Port        string
	DBConn      *gorm.DB
	ServerReady chan bool
}

// Start start http server
func (s *Server) Start() {
	itemHandler := newHandler(s)
	router := gin.Default()

	router.Use(newServerMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"pong": "ok",
		})
		return
	})

	api := router.Group("/v1/items")
	{
		api.GET("/:id", itemHandler.GetItemByID)
		api.POST("/", itemHandler.CreateItem)
	}

	if err := router.Run(":" + s.Port); err != nil {
		logrus.Errorf(err.Error())
		logrus.Infof("shutting down the server")
	}
}

func newHandler(s *Server) handler.ItemHTTPHandler {
	itemRepository := repository.NewItemRepository(s.DBConn)
	itemService := service.NewItemService(itemRepository)
	itemHandler := handler.NewItemHTTPHandler(itemService)

	return itemHandler
}
