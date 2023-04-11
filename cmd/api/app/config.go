package app

import (
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/market/item"
	"github.com/osalomon89/test-crud-api/internal/market/user"
	"github.com/osalomon89/test-crud-api/internal/market/user/token"
	"github.com/osalomon89/test-crud-api/internal/platform/environment"
	"github.com/osalomon89/test-crud-api/internal/platform/jwt"
	"github.com/osalomon89/test-crud-api/internal/platform/mysql"
)

type Dependencies struct {
	ItemRepository item.Repository
	UserRepository user.Repository
	TokenService   token.Service
}

func BuildDependencies(env environment.Environment) (*Dependencies, error) {
	switch env {
	case environment.Development:
		mysqlDb, err := mysql.GetConnectionDB()
		if err != nil {
			return nil, fmt.Errorf("error connecting to DB: %w", err)
		}

		// infra adapters
		itemRepo, err := item.NewItemRepository(mysqlDb)
		if err != nil {
			return nil, fmt.Errorf("error connecting to DB: %w", err)
		}

		userRepo, err := user.NewUserRepository(mysqlDb)
		if err != nil {
			return nil, fmt.Errorf("error connecting to DB: %w", err)
		}

		tokenService := token.NewTokenGenerator(jwt.New())

		return &Dependencies{
			ItemRepository: itemRepo,
			UserRepository: userRepo,
			TokenService:   tokenService,
		}, nil
	}

	return nil, nil
}