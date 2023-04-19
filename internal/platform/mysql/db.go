package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserDAO struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

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
	Photos      []string
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

type MySQLDB struct {
	conn *sqlx.DB
}

func NewMySQLDB(conn *sqlx.DB) (*MySQLDB, error) {
	if conn == nil {
		return nil, fmt.Errorf("mysql connection cannot be nil")
	}

	return &MySQLDB{conn: conn}, nil
}

func (repo *MySQLDB) SaveItem(ctx context.Context, item *ItemDAO) error {
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
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if len(item.Photos) > 0 {
		if err := repo.savePhotos(tx, uint(id), item.Photos); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	item.ID = uint(id)
	item.CreatedAt = createdAt
	item.UpdatedAt = createdAt

	return nil
}

func (repo *MySQLDB) GetItemByID(ctx context.Context, id uint) (*ItemDAO, error) {
	item := new(ItemDAO)
	err := repo.conn.Get(item, "SELECT * FROM items WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	photos := []PhotoDAO{}
	err = repo.conn.Select(&photos, "SELECT * FROM photos WHERE item_id=?", id)
	if err != nil {
		return nil, err
	}

	for _, photo := range photos {
		item.Photos = append(item.Photos, photo.Path)
	}

	return item, nil
}

func (repo *MySQLDB) savePhotos(tx *sql.Tx, id uint, photos []string) error {
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

func (repo *MySQLDB) SaveUser(ctx context.Context, user *UserDAO) error {
	createdAt := time.Now()

	result, err := repo.conn.Exec(`INSERT INTO users (email, password, created_at, updated_at) 
	VALUES(?,?,?,?)`, user.Email, user.Password, createdAt, createdAt)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	user.CreatedAt = createdAt
	user.UpdatedAt = createdAt

	return nil
}

func (repo *MySQLDB) GetUserByEmail(ctx context.Context, email string) (*UserDAO, error) {
	user := new(UserDAO)

	err := repo.conn.Get(user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *MySQLDB) SaveItemByUser(ctx context.Context, itemID, userID uint) error {
	createdAt := time.Now()

	_, err := repo.conn.Exec(`INSERT INTO user_items (user_id, item_id, created_at, updated_at) 
	VALUES(?,?,?,?)`, itemID, userID, createdAt, createdAt)

	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	return nil
}

func (repo *MySQLDB) GetItemsByUser(ctx context.Context, userID uint) ([]ItemDAO, error) {
	var items []ItemDAO

	err := repo.conn.Select(&items, `SELECT * FROM items 
	INNER JOIN user_items ON items.id = user_items.item_id 
	WHERE user_items.user_id=?`, userID)
	if err != nil {
		return items, err
	}

	return items, nil
}
