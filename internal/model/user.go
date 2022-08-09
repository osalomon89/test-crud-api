package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique"`
	Items []Item `gorm:"many2many:user_items;"`
}

type UserItem struct {
	gorm.Model
	UserID uint `gorm:"primaryKey"`
	ItemID uint `gorm:"primaryKey"`
}
