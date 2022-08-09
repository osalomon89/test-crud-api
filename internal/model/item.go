package model

import "gorm.io/gorm"

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
	gorm.Model
	Code        string `gorm:"unique"`
	Title       string
	Description string
	Price       int
	Stock       int
	ItemType    string
	Leader      bool
	LeaderLevel string
	Status      string
	Photos      []Photo
}

type Photo struct {
	gorm.Model
	Path   string
	ItemID uint
}
