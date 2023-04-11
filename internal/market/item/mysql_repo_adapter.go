package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type ItemDAO struct {
	ID          uint
	Code        string
	Title       string
	Description string
	Price       int
	Stock       int
	ItemType    string `db:"item_type"`
	Leader      bool
	LeaderLevel string `db:"leader_level"`
	Status      string
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PhotoDAO struct {
	ID        uint
	Path      string
	ItemID    uint      `db:"item_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type itemRepository struct {
	conn *sqlx.DB
}

func NewItemRepository(conn *sqlx.DB) (Repository, error) {
	if conn == nil {
		return nil, fmt.Errorf("mysql connection cannot be nil")
	}

	return &itemRepository{conn: conn}, nil
}

func (repo *itemRepository) Save(ctx context.Context, item *Item) error {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering ItemRepository. CreateItem()")

	tx, err := repo.conn.Begin()
	if err != nil {
		return fmt.Errorf("transaction initialization error: %w", err)
	}

	createdAt := time.Now()
	result, err := tx.Exec(`INSERT INTO items 
		(code, title, description, price, stock, item_type, leader, leader_level, status, created_at, updated_at) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?)`, item.Code, item.Title, item.Description, item.Price, item.Stock,
		item.ItemType, item.Leader, item.LeaderLevel, item.Status, createdAt, createdAt)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ItemError{
				Message: "The item code must be unique",
			}
		}

		return fmt.Errorf("error saving item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	if len(item.Photos) > 0 {
		if err := repo.savePhotos(tx, uint(id), item.Photos); err != nil {
			return fmt.Errorf("error saving photos: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	item.ID = uint(id)
	item.CreatedAt = createdAt
	item.UpdatedAt = createdAt

	return nil
}

func (repo *itemRepository) GetItemByID(ctx context.Context, id uint) (*Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering ItemRepository. GetItemByID()")

	item := new(ItemDAO)
	err := repo.conn.Get(item, "SELECT * FROM items WHERE id=?", id)
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

	photos := []PhotoDAO{}
	err = repo.conn.Select(&photos, "SELECT * FROM photos WHERE item_id=?", id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			break
		default:
			return nil, fmt.Errorf("error getting photos: %w", err)
		}
	}

	return item.toItemDomain(photos), nil
}

func (item *ItemDAO) toItemDomain(photos []PhotoDAO) *Item {
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
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	for _, photo := range photos {
		itemModel.Photos = append(itemModel.Photos, photo.Path)
	}

	return &itemModel
}

func (repo *itemRepository) savePhotos(tx *sql.Tx, id uint, photos []string) error {
	createdAt := time.Now()
	valueStrings := make([]string, 0, len(photos))
	valueArgs := make([]interface{}, 0, len(photos)*4)

	for _, path := range photos {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, path)
		valueArgs = append(valueArgs, id)
		valueArgs = append(valueArgs, createdAt)
		valueArgs = append(valueArgs, createdAt)
	}

	stmt := fmt.Sprintf(`INSERT INTO photos (path, item_id, created_at, updated_at) VALUES %s`,
		strings.Join(valueStrings, ","))

	_, err := tx.Exec(stmt, valueArgs...)

	return err
}
