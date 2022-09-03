package services

import (
	"context"
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/core/domain"
	"github.com/osalomon89/test-crud-api/internal/core/ports"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type itemService struct {
	itemRepository ports.ItemRepository
}

func NewItemService(itemRepository ports.ItemRepository) (ports.ItemService, error) {
	if itemRepository == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}

	return &itemService{itemRepository: itemRepository}, nil
}

func (svc *itemService) CreateItem(ctx context.Context,
	item domain.Item) (*domain.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(svc, nil, "Entering ItemService. CreateItem()")

	item.SetStatus()

	if err := validateItemModel(&item); err != nil {
		return nil, err
	}

	err := svc.itemRepository.SaveItem(ctx, &item)
	if err != nil {
		return nil, fmt.Errorf("error in repository: %w", err)
	}

	return &item, nil
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

func validateItemModel(item *domain.Item) error {
	if item.ItemType == domain.ItemTypeSeller {
		if item.Leader {
			if !isAValidLeaderLevel(item.LeaderLevel) {
				return domain.ItemError{
					Message: fmt.Sprintf("Error in params validation. Leader level is not valid: %s", item.LeaderLevel),
				}
			}
		} else {
			item.LeaderLevel = ""
		}
	}

	if len(item.Photos) == 0 {
		return domain.ItemError{
			Message: "Error in params validation: photos array len can not be null",
		}
	}

	return nil
}

func isAValidLeaderLevel(leaderLever string) bool {
	if leaderLever == domain.LeaderLevelBasic ||
		leaderLever == domain.LeaderLevelGold ||
		leaderLever == domain.LeaderLevelPlatinum {
		return true
	}

	return false
}
