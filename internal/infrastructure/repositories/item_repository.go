package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/application/ports"
	"github.com/osalomon89/test-crud-api/internal/domain"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Code        string `gorm:"unique"`
	Title       string
	Description string
	Price       int
	Stock       int
	ItemType    string
	Leader      bool
	LeaderLevel string
	Status      string
	Photos      []Photo
}

type Photo struct {
	gorm.Model
	Path   string
	ItemID uint
}

type itemRepository struct {
	conn *gorm.DB
}

func NewItemRepository(conn *gorm.DB) (ports.ItemRepository, error) {
	if conn == nil {
		return nil, fmt.Errorf("mysql connection cannot be nil")
	}

	return &itemRepository{conn: conn}, nil
}

func (r *itemRepository) CreateItem(ctx context.Context, item *domain.Item) error {
	logger := marketcontext.Logger(ctx)
	logger.Debug(r, nil, "Entering ItemRepository. CreateItem()")

	itemModel := r.marshalItem(item)

	if err := r.conn.Create(itemModel).Error; err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	item.ID = itemModel.ID
	item.CreatedAt = itemModel.CreatedAt
	item.UpdatedAt = itemModel.UpdatedAt

	return nil
}

func (r *itemRepository) GetItemByID(ctx context.Context, id uint) (*domain.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(r, nil, "Entering ItemRepository. GetItemByID()")

	item := new(Item)
	if err := r.conn.First(item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ResourceNotFoundError{
				Message: "Item not found",
			}
		}

		return nil, fmt.Errorf("error getting items %v: %w", id, err)
	}

	var photos []Photo
	err := r.conn.Where("item_id = ?", id).Find(&photos).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ResourceNotFoundError{
				Message: "Photos not found",
			}
		}

		return nil, fmt.Errorf("error getting photos %v: %w", id, err)
	}

	item.Photos = photos

	return r.unmarshalItem(item), nil
}

func (r *itemRepository) marshalItem(item *domain.Item) *Item {
	itemModel := Item{
		Code:        item.Code,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		ItemType:    item.ItemType,
		Leader:      item.Leader,
		LeaderLevel: item.LeaderLevel,
		Status:      item.Status,
	}

	for _, photo := range item.Photos {
		itemModel.Photos = append(itemModel.Photos, Photo{Path: photo.Path})
	}

	return &itemModel
}

func (r *itemRepository) unmarshalItem(item *Item) *domain.Item {
	itemModel := domain.Item{
		ID:          item.ID,
		Code:        item.Code,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		ItemType:    item.ItemType,
		Leader:      item.Leader,
		LeaderLevel: item.LeaderLevel,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	for _, photo := range item.Photos {
		itemModel.Photos = append(itemModel.Photos, domain.Photo{Path: photo.Path})
	}

	return &itemModel
}
