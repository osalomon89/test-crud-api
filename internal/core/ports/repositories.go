package ports

import (
	"context"

	"github.com/osalomon89/test-crud-api/internal/core/domain"
)

//go:generate mockgen github.com/osalomon89/test-crud-api/internal/core/ports ItemRepository -destination=../test/mocks/item_repository_mock.go -package=mocks
type ItemRepository interface {
	SaveItem(ctx context.Context, a *domain.Item) error
	GetItemByID(ctx context.Context, id uint) (*domain.Item, error)
}
