package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/test-crud-api/internal/market/item"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type userDAO struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) (Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("database can not be nil")
	}

	return &userRepository{db: db}, nil
}

func (repo *userRepository) Save(ctx context.Context, user *User) error {
	createdAt := time.Now()

	result, err := repo.db.Exec(`INSERT INTO users (email, password, created_at, updated_at) 
	VALUES(?,?,?,?)`, user.Email, user.Password, createdAt, createdAt)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return UserError{
				Message: "The email must be unique",
			}
		}

		return fmt.Errorf("error saving user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	user.ID = uint(id)
	user.CreatedAt = createdAt
	user.UpdatedAt = createdAt

	return nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := new(userDAO)

	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering UserRepository. GetByEmail()")

	err := repo.db.Get(user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ResourceNotFoundError{
				Message: "user not found",
			}
		default:
			return nil, fmt.Errorf("error getting user: %w", err)
		}
	}

	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (repo *userRepository) SaveItem(ctx context.Context, itemID, userID uint) error {
	createdAt := time.Now()

	_, err := repo.db.Exec(`INSERT INTO user_items (user_id, item_id, created_at, updated_at) 
	VALUES(?,?,?,?)`, itemID, userID, createdAt, createdAt)

	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	return nil
}

func (repo *userRepository) GetItems(ctx context.Context, userID uint) (*Items, error) {
	var items []item.Item

	err := repo.db.Select(&items, `SELECT * FROM items 
	INNER JOIN user_items ON items.id = user_items.item_id 
	WHERE user_items.user_id=?`, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ResourceNotFoundError{
				Message: "user not found",
			}
		default:
			return nil, fmt.Errorf("error getting user: %w", err)
		}
	}

	return &Items{
		Pagination: Pagination{
			Page:       1,
			PageSize:   10,
			TotalPages: len(items),
			Total:      len(items),
		},
		Data: items,
	}, nil
}
