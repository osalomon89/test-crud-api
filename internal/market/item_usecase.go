package market

import (
	"context"
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/market/item"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type CRUDUseCase interface {
	Create(ctx context.Context, item item.Item) (*item.Item, error)
	Read(ctx context.Context, id uint) (*item.Item, error)
}

type useCaseCRUD struct {
	repo item.Repository
}

func NewCRUDUseCase(repo item.Repository) *useCaseCRUD {
	return &useCaseCRUD{repo: repo}
}

func (c useCaseCRUD) Create(ctx context.Context, item item.Item) (*item.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(c, nil, "Entering ItemService. CreateItem()")

	item.SetStatus()

	if err := validateItemModel(&item); err != nil {
		return nil, err
	}

	if err := c.repo.Save(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (c useCaseCRUD) Read(ctx context.Context, id uint) (*item.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(c, nil, "Entering ItemService. GetItemByID()")

	item, err := c.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error in repository: %w", err)
	}

	return item, nil
}

func validateItemModel(itemObj *item.Item) error {
	if itemObj.ItemType == item.ItemTypeSeller {
		if itemObj.Leader {
			if !isAValidLeaderLevel(itemObj.LeaderLevel) {
				return item.ItemError{
					Message: fmt.Sprintf("Error in params validation. Leader level is not valid: %s", itemObj.LeaderLevel),
				}
			}
		} else {
			itemObj.LeaderLevel = ""
		}
	}

	if len(itemObj.Photos) == 0 {
		return item.ItemError{
			Message: "Error in params validation: photos array len can not be null",
		}
	}

	return nil
}

func isAValidLeaderLevel(leaderLever string) bool {
	if leaderLever == item.LeaderLevelBasic ||
		leaderLever == item.LeaderLevelGold ||
		leaderLever == item.LeaderLevelPlatinum {
		return true
	}

	return false
}
