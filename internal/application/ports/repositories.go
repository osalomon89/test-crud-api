package ports

import (
	"context"

	"github.com/osalomon89/test-crud-api/internal/domain"
)

//go:generate mockgen -destination=../utils/test/mocks/item_repository_mock.go -package=mocks -source=./item_repository.go
type ItemRepository interface {
	SaveItem(ctx context.Context, a *domain.Item) error
	GetItemByID(ctx context.Context, id uint) (*domain.Item, error)
}
