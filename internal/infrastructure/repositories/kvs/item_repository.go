package kvs

import (
	"context"
	"fmt"
	"time"

	"github.com/mercadolibre/go-meli-toolkit/gokvsclient"
	"github.com/osalomon89/test-crud-api/internal/application/ports"
	"github.com/osalomon89/test-crud-api/internal/domain"
)

type Item struct {
	Code        string
	Title       string
	Description string
	Price       int
	Stock       int
	ItemType    string
	Leader      bool
	LeaderLevel string
	Status      string
	Photos      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type itemRepository struct {
	client gokvsclient.Client
}

func NewItemRepository(client gokvsclient.Client) (ports.ItemRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("gokvsclient cannot be nil")
	}

	return &itemRepository{client: client}, nil
}

func (repo *itemRepository) SaveItem(ctx context.Context, item *domain.Item) error {
	itemExist, err := repo.itemExist(item.Code)
	if err != nil {
		return err
	}

	if itemExist != nil {
		return domain.ItemError{
			Message: "The item code must be unique",
		}
	}

	createdAt := time.Now()
	itemDTO := new(Item)

	itemDTO.Code = item.Code
	itemDTO.Title = item.Title
	itemDTO.Description = item.Description
	itemDTO.Price = item.Price
	itemDTO.Stock = item.Stock
	itemDTO.ItemType = item.ItemType
	itemDTO.Leader = item.Leader
	itemDTO.LeaderLevel = item.LeaderLevel
	itemDTO.Status = item.Status
	itemDTO.CreatedAt = createdAt
	itemDTO.UpdatedAt = createdAt

	for _, photo := range item.Photos {
		itemDTO.Photos = append(itemDTO.Photos, photo.Path)
	}

	kvsItem := gokvsclient.MakeItem(itemDTO.Code, itemDTO)
	if err := repo.client.Save(kvsItem); err != nil {
		return fmt.Errorf("error saving item in KVS: %w", err)
	}

	item.CreatedAt = createdAt
	item.UpdatedAt = createdAt

	return nil
}

func (repo *itemRepository) itemExist(code string) (gokvsclient.Item, error) {
	kvsItem, err := repo.client.Get(code)
	if err != nil {
		return nil, fmt.Errorf("error getting item: %w", err)
	}

	return kvsItem, nil
}

func (repo *itemRepository) GetItemByID(ctx context.Context, id uint) (*domain.Item, error) {
	kvsItem, err := repo.client.Get("SAM27324355")
	if err != nil {
		return nil, fmt.Errorf("error getting item: %w", err)
	}

	item := new(Item)
	if err := kvsItem.GetValue(item); err != nil {
		return nil, fmt.Errorf("error unmarshaling item: %w", err)
	}

	return repo.createItem(item), nil
}

func (repo *itemRepository) createItem(itemDTO *Item) *domain.Item {
	item := new(domain.Item)

	item.Code = itemDTO.Code
	item.Title = itemDTO.Title
	item.Description = itemDTO.Description
	item.Price = itemDTO.Price
	item.Stock = itemDTO.Stock
	item.ItemType = itemDTO.ItemType
	item.Leader = itemDTO.Leader
	item.LeaderLevel = itemDTO.LeaderLevel
	item.Status = itemDTO.Status
	item.CreatedAt = itemDTO.CreatedAt
	item.UpdatedAt = itemDTO.UpdatedAt

	for k, photo := range item.Photos {
		item.Photos = append(item.Photos, domain.Photo{
			ID:   uint(k),
			Path: photo.Path,
		})
	}

	return item
}
