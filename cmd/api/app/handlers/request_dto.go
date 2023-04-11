package handlers

import (
	"github.com/osalomon89/test-crud-api/internal/market/item"
)

type itemRequest struct {
	Code        string   `json:"code" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       int      `json:"price" binding:"required"`
	Stock       int      `json:"stock" binding:"required"`
	ItemType    string   `json:"itemType" binding:"required"`
	Leader      bool     `json:"leader"`
	LeaderLevel string   `json:"leaderLevel"`
	Photos      []string `json:"photos" binding:"required"`
}

func (itemBody itemRequest) toItemDomain() item.Item {
	return item.Item{
		Code:        itemBody.Code,
		Title:       itemBody.Title,
		Description: itemBody.Description,
		Price:       itemBody.Price,
		Stock:       itemBody.Stock,
		ItemType:    itemBody.ItemType,
		Leader:      itemBody.Leader,
		LeaderLevel: itemBody.LeaderLevel,
		Photos:      itemBody.Photos,
	}
}
