package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/osalomon89/test-crud-api/internal/market/item"
	mysqldb "github.com/osalomon89/test-crud-api/internal/platform/mysql"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type userRepository struct {
	db *mysqldb.MySQLDB
}

func NewUserRepository(db *mysqldb.MySQLDB) (Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("database can not be nil")
	}

	return &userRepository{db: db}, nil
}

func (repo *userRepository) Save(ctx context.Context, user *User) error {
	userDao := mysqldb.UserDAO{
		Email:    user.Email,
		Password: user.Password,
	}

	if err := repo.db.SaveUser(ctx, &userDao); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return UserError{
				Message: "The email must be unique",
			}
		}

		return fmt.Errorf("error saving user: %w", err)
	}

	user.ID = userDao.ID
	user.CreatedAt = userDao.CreatedAt
	user.UpdatedAt = userDao.UpdatedAt

	return nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering UserRepository. GetByEmail()")

	userDao, err := repo.db.GetUserByEmail(ctx, email)
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
		ID:        userDao.ID,
		Email:     userDao.Email,
		Password:  userDao.Password,
		CreatedAt: userDao.CreatedAt,
		UpdatedAt: userDao.UpdatedAt,
	}, nil
}

func (repo *userRepository) SaveItem(ctx context.Context, itemID, userID uint) error {
	logger := marketcontext.Logger(ctx)
	logger.Debug(repo, nil, "Entering UserRepository. SaveItem()")

	if err := repo.db.SaveItemByUser(ctx, itemID, userID); err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	return nil
}

func (repo *userRepository) GetItems(ctx context.Context, userID uint) (*Items, error) {
	var items []item.Item

	itemsDao, err := repo.db.GetItemsByUser(ctx, userID)
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

	for _, it := range itemsDao {
		items = append(items, *item.ToItemDomain(&it))
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
