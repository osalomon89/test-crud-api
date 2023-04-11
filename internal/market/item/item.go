package item

import (
	"context"
	"time"
)

const (
	ItemTypeOwn         = "OWN"
	ItemTypeSeller      = "SELLER"
	StatusAll           = "ALL"
	StatusActive        = "ACTIVE"
	StatusInactive      = "INACTIVE"
	LeaderLevelBasic    = "BASIC"
	LeaderLevelGold     = "GOLD"
	LeaderLevelPlatinum = "PLATINUM"
)

type Item struct {
	ID          uint
	Code        string
	Title       string
	Description string
	Price       int
	Stock       int
	ItemType    string
	Leader      bool
	LeaderLevel string
	Status      string
	Photos      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repository interface {
	Save(ctx context.Context, a *Item) error
	GetItemByID(ctx context.Context, id uint) (*Item, error)
}

func (item *Item) SetStatus() {
	if item.Stock > 0 {
		item.Status = StatusActive
		return
	}

	item.Status = StatusInactive
}
