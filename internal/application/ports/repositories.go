package ports

import (
	"context"

	"github.com/osalomon89/test-crud-api/internal/domain"
)

//go:generate mockgen -source=./repositories.go -destination=../test/mocks/item_repository_mock.go -package=mocks
type ItemRepository interface {
	SaveItem(ctx context.Context, a *domain.Item) error
	GetItemByID(ctx context.Context, id uint) (*domain.Item, error)
}
