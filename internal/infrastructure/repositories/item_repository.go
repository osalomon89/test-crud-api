package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/test-crud-api/internal/application/ports"
	"github.com/osalomon89/test-crud-api/internal/domain"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type Item struct {
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
	Photos      []Photo
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Photo struct {
	ID        uint
	Path      string
	ItemID    uint      `db:"item_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type itemRepository struct {
	conn *sqlx.DB
}

func NewItemRepository(conn *sqlx.DB) (ports.ItemRepository, error) {
	if conn == nil {
		return nil, fmt.Errorf("mysql connection cannot be nil")
	}

	return &itemRepository{conn: conn}, nil
}

func (repo *itemRepository) SaveItem(ctx context.Context, item *domain.Item) error {
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
			return domain.ItemError{
				Message: "The item code must be unique",
			}
		}

		return fmt.Errorf("error saving item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	err = repo.savePhotos(tx, uint(id), item.Photos)
	if err != nil {
		return fmt.Errorf("error saving photos: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	item.ID = uint(id)
	item.CreatedAt = createdAt
	item.UpdatedAt = createdAt

	return nil
}

func (repo *itemRepository) GetItemByID(ctx context.Context, id uint) (*domain.Item, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering ItemRepository. GetItemByID()")

	item := new(Item)
	err := repo.conn.Get(item, "SELECT * FROM items WHERE id=?", id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, domain.ResourceNotFoundError{
				Message: "Item not found",
			}
		default:
			return nil, fmt.Errorf("error getting items: %w", err)
		}
	}

	return repo.unmarshalItem(item), nil
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

func (repo *itemRepository) savePhotos(tx *sql.Tx, id uint, photos []domain.Photo) error {
	createdAt := time.Now()
	valueStrings := make([]string, 0, len(photos))
	valueArgs := make([]interface{}, 0, len(photos)*4)

	for _, photo := range photos {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, photo.Path)
		valueArgs = append(valueArgs, id)
		valueArgs = append(valueArgs, createdAt)
		valueArgs = append(valueArgs, createdAt)
	}

	stmt := fmt.Sprintf(`INSERT INTO photos (path, item_id, created_at, updated_at) VALUES %s`,
		strings.Join(valueStrings, ","))

	_, err := tx.Exec(stmt, valueArgs...)

	return err
}
