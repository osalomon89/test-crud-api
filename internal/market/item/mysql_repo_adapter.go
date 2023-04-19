package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	mysqldb "github.com/osalomon89/test-crud-api/internal/platform/mysql"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type itemRepository struct {
	conn *mysqldb.MySQLDB
}

func NewItemRepository(conn *mysqldb.MySQLDB) (Repository, error) {
	if conn == nil {
		return nil, fmt.Errorf("mysql connection cannot be nil")
	}

	return &itemRepository{conn: conn}, nil
}

func (repo *itemRepository) Save(ctx context.Context, item *Item) error {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering ItemRepository. CreateItem()")

	itemDao := mysqldb.ItemDAO{
		Code:        item.Code,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		ItemType:    item.ItemType,
		Leader:      item.Leader,
		LeaderLevel: item.LeaderLevel,
		Status:      item.Status,
		Photos:      item.Photos,
	}

	err := repo.conn.SaveItem(ctx, &itemDao)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ItemError{
				Message: "The item code must be unique",
			}
		}

		return fmt.Errorf("error saving item: %w", err)
	}

	item.ID = itemDao.ID
	item.CreatedAt = itemDao.CreatedAt
	item.UpdatedAt = itemDao.UpdatedAt

	return nil
}

func (repo *itemRepository) GetItemByID(ctx context.Context, id uint) (*Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering ItemRepository. GetItemByID()")

	itemDao, err := repo.conn.GetItemByID(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ResourceNotFoundError{
				Message: "Item not found",
			}
		default:
			return nil, fmt.Errorf("error getting items: %w", err)
		}
	}

	return ToItemDomain(itemDao), nil
}

func ToItemDomain(item *mysqldb.ItemDAO) *Item {
	itemModel := Item{
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
		Photos:      item.Photos,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	return &itemModel
}
