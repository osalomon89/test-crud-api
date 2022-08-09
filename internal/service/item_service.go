package service

import (
	"context"
	"fmt"

	itemerrors "github.com/osalomon89/test-crud-api/internal/errors"
	"github.com/osalomon89/test-crud-api/internal/model"
	httpmodel "github.com/osalomon89/test-crud-api/internal/model/http"
	"github.com/osalomon89/test-crud-api/internal/repository"
	"github.com/osalomon89/test-crud-api/internal/server/httpcontext"
)

type ItemService interface {
	CreateItem(ctx context.Context, itemBody httpmodel.CreateItemBody) (*model.Item, error)
	GetItemByID(ctx context.Context, itemID uint) (*model.Item, error)
}

type itemService struct {
	itemRepository repository.ItemRepository
}

func NewItemService(itemRepository repository.ItemRepository) ItemService {
	return &itemService{itemRepository: itemRepository}
}

func (srv *itemService) CreateItem(ctx context.Context,
	itemBody httpmodel.CreateItemBody) (*model.Item, error) {
	logger := httpcontext.GetLogger(ctx)
	logger.Info("Entering ItemService. CreateItem()")

	item, err := validateItemModel(createItemModel(itemBody))
	if err != nil {
		logger.Info("ItemService: error in service: %w", err)
		return nil, err
	}

	err = srv.itemRepository.CreateItem(ctx, item)
	if err != nil {
		logger.Info("ItemService: error in repository: %w", err)
		return nil, fmt.Errorf("error in repository: %w", err)
	}

	return item, nil
}

func (srv *itemService) GetItemByID(ctx context.Context, itemID uint) (*model.Item, error) {
	logger := httpcontext.GetLogger(ctx)
	logger.Info("Entering ItemService. GetItemByID()")

	item, err := srv.itemRepository.GetItemByID(ctx, itemID)
	if err != nil {
		logger.Info("ItemService: error in repository: %w", err)
		return nil, fmt.Errorf("error in repository: %w", err)
	}

	return item, nil
}

func createItemModel(itemBody httpmodel.CreateItemBody) *model.Item {
	item := new(model.Item)

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

func getModelPhotos(photosList []string) []model.Photo {
	var photos []model.Photo

	for _, path := range photosList {
		photos = append(photos, model.Photo{Path: path})
	}

	return photos
}

func validateItemModel(item *model.Item) (*model.Item, error) {
	if item.ItemType == model.ItemTypeSeller {
		if item.Leader {
			if !isAValidLeaderLevel(item.LeaderLevel) {
				return nil, itemerrors.ItemError{
					Message: fmt.Sprintf("Error in params validation. Leader level is not valid: %s", item.LeaderLevel),
				}
			}
		} else {
			item.LeaderLevel = ""
		}
	}

	if len(item.Photos) == 0 {
		return nil, itemerrors.ItemError{
			Message: "Error in params validation: photos array len can not be null",
		}
	}

	return item, nil
}

func isAValidLeaderLevel(leaderLever string) bool {
	if leaderLever == model.LeaderLevelBasic ||
		leaderLever == model.LeaderLevelGold ||
		leaderLever == model.LeaderLevelPlatinum {
		return true
	}

	return false
}
