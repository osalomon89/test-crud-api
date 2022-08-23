package ports

import (
	"context"

	"github.com/osalomon89/test-crud-api/internal/domain"
)

type ItemService interface {
	CreateItem(ctx context.Context, itemBody domain.CreateItemBody) (*domain.Item, error)
	GetItemByID(ctx context.Context, itemID uint) (*domain.Item, error)
}
