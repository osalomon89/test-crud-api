package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/osalomon89/test-crud-api/cmd/api/app/handlers"
	"github.com/osalomon89/test-crud-api/cmd/api/app/middleware"
	"github.com/osalomon89/test-crud-api/internal/market"
)

func Build(dep *Dependencies) *gin.Engine {
	r := gin.Default()

	// use cases
	crudUseCase := market.NewCRUDUseCase(dep.ItemRepository)
	userUseCase := market.NewUserUsecase(dep.UserRepository,
		dep.ItemRepository, dep.TokenService)

	// middleware
	md := middleware.NewAuthMiddleware(dep.TokenService)

	// controller adapters
	crudHandler, err := handlers.NewItemHandler(crudUseCase)
	if err != nil {
		log.Fatal(err)
	}

	userHandler, err := handlers.NewUserHandler(userUseCase)
	if err != nil {
		log.Fatal(err)
	}

	basePath := "/api/v1/market"
	publicRouter := r.Group(basePath)

	publicRouter.GET("/items/{id}", crudHandler.GetItemByID)
	publicRouter.POST("/items/", crudHandler.CreateItem)
	publicRouter.POST("/users", userHandler.Register)
	publicRouter.POST("/users/login", userHandler.Login)

	protectedRouter := r.Group(fmt.Sprintf("%s/users/me", basePath))
	// Middleware to verify AccessToken
	protectedRouter.Use(md.TokenAuthMiddleware())

	protectedRouter.Use(md.TokenAuthMiddleware())

	protectedRouter.GET("/favorites", userHandler.MarkItemAsFavorite)
	//protectedRouter.POST("/favorites", userHandler.SaveProperty)

	return r
}
