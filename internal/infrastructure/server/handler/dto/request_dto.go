package dto

import "github.com/osalomon89/test-crud-api/internal/core/domain"

type ItemBody struct {
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

func (itemBody ItemBody) ToItemDomain() domain.Item {
	var photos []domain.Photo

	for _, path := range itemBody.Photos {
		photos = append(photos, domain.Photo{Path: path})
	}

	return domain.Item{
		Code:        itemBody.Code,
		Title:       itemBody.Title,
		Description: itemBody.Description,
		Price:       itemBody.Price,
		Stock:       itemBody.Stock,
		ItemType:    itemBody.ItemType,
		Leader:      itemBody.Leader,
		LeaderLevel: itemBody.LeaderLevel,
		Photos:      photos,
	}
}
