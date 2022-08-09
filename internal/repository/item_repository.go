package repository

import (
	"context"
	"errors"
	"fmt"

	itemerrors "github.com/osalomon89/test-crud-api/internal/errors"
	"github.com/osalomon89/test-crud-api/internal/model"
	"github.com/osalomon89/test-crud-api/internal/server/httpcontext"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../utils/test/mocks/item_repository_mock.go -package=mocks -source=./item_repository.go
type ItemRepository interface {
	CreateItem(ctx context.Context, a *model.Item) error
	GetItemByID(ctx context.Context, id uint) (*model.Item, error)
}

// itemRepository implements ItemRepository interface
type itemRepository struct {
	conn *gorm.DB
}

func NewItemRepository(conn *gorm.DB) ItemRepository {
	return &itemRepository{conn: conn}
}

func (r *itemRepository) CreateItem(ctx context.Context, item *model.Item) error {
	logger := httpcontext.GetLogger(ctx)
	logger.Info("Entering ItemRepository. CreateItem()")

	if err := r.conn.Create(item).Error; err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	return nil
}

func (r *itemRepository) GetItemByID(ctx context.Context, id uint) (*model.Item, error) {
	logger := httpcontext.GetLogger(ctx)
	logger.Info("Entering ItemRepository. GetItemByID()")

	item := new(model.Item)
	if err := r.conn.First(item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, itemerrors.ResourceNotFoundError{
				Message: "Item not found",
			}
		}

		return nil, fmt.Errorf("error getting items %v: %w", id, err)
	}

	var photos []model.Photo
	err := r.conn.Where("item_id = ?", id).Find(&photos).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, itemerrors.ResourceNotFoundError{
				Message: "Photos not found",
			}
		}

		return nil, fmt.Errorf("error getting photos %v: %w", id, err)
	}

	item.Photos = photos

	return item, nil
}
