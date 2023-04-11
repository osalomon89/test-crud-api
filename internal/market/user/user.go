package user

import (
	"context"
	"time"

	"github.com/osalomon89/test-crud-api/internal/market/item"
)

type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Items struct {
	Pagination
	Data []item.Item
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	Total      int `json:"total"`
}

type Repository interface {
	Save(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	SaveItem(ctx context.Context, itemID, userID uint) error
	GetItems(ctx context.Context, userID uint) (*Items, error)
}
