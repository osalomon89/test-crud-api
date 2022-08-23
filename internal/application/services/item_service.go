package services

import (
	"context"
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/application/ports"
	"github.com/osalomon89/test-crud-api/internal/domain"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type itemService struct {
	itemRepository ports.ItemRepository
}

func NewItemService(itemRepository ports.ItemRepository) ports.ItemService {
	return &itemService{itemRepository: itemRepository}
}

func (svc *itemService) CreateItem(ctx context.Context,
	itemBody domain.CreateItemBody) (*domain.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(svc, nil, "Entering ItemService. CreateItem()")

	item, err := validateItemModel(createItemModel(itemBody))
	if err != nil {
		return nil, err
	}

	err = svc.itemRepository.CreateItem(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("error in repository: %w", err)
	}

	return item, nil
}

func (svc *itemService) GetItemByID(ctx context.Context, itemID uint) (*domain.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(svc, nil, "Entering ItemService. GetItemByID()")

	item, err := svc.itemRepository.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("error in repository: %w", err)
	}

	return item, nil
}

func createItemModel(itemBody domain.CreateItemBody) *domain.Item {
	item := new(domain.Item)

	item.Code = itemBody.Code
	item.Title = itemBody.Title
	item.Description = itemBody.Description
	item.Price = itemBody.Price
	item.Stock = itemBody.Stock
	item.ItemType = itemBody.ItemType
	item.Leader = itemBody.Leader
	item.LeaderLevel = itemBody.LeaderLevel
	item.Photos = getModelPhotos(itemBody.Photos)
	item.Status = itemBody.GetStatus()

	return item
}

func getModelPhotos(photosList []string) []domain.Photo {
	var photos []domain.Photo

	for _, path := range photosList {
		photos = append(photos, domain.Photo{Path: path})
	}

	return photos
}

func validateItemModel(item *domain.Item) (*domain.Item, error) {
	if item.ItemType == domain.ItemTypeSeller {
		if item.Leader {
			if !isAValidLeaderLevel(item.LeaderLevel) {
				return nil, domain.ItemError{
					Message: fmt.Sprintf("Error in params validation. Leader level is not valid: %s", item.LeaderLevel),
				}
			}
		} else {
			item.LeaderLevel = ""
		}
	}

	if len(item.Photos) == 0 {
		return nil, domain.ItemError{
			Message: "Error in params validation: photos array len can not be null",
		}
	}

	return item, nil
}

func isAValidLeaderLevel(leaderLever string) bool {
	if leaderLever == domain.LeaderLevelBasic ||
		leaderLever == domain.LeaderLevelGold ||
		leaderLever == domain.LeaderLevelPlatinum {
		return true
	}

	return false
}
