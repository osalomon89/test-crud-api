package domain

import "time"

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
	Photos      []Photo
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Photo struct {
	ID        uint
	Path      string
	ItemID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (item *Item) SetStatus() {
	if item.Stock > 0 {
		item.Status = StatusActive
		return
	}

	item.Status = StatusInactive
}
