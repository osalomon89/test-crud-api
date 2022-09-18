package handlers

import (
	"time"

	"github.com/osalomon89/test-crud-api/internal/core/domain"
)

type Response struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    *ItemResponse `json:"data"`
}

type ItemResponse struct {
	ID          uint     `json:"id"`
	Code        string   `json:"code"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Stock       int      `json:"stock"`
	ItemType    string   `json:"itemType"`
	Leader      bool     `json:"leader"`
	LeaderLevel string   `json:"leaderLevel"`
	Status      string   `json:"status"`
	Photos      []string `json:"photos"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (itemResponse *ItemResponse) fromItemDomain(item *domain.Item) {
	var photos []string
	for _, photo := range item.Photos {
		photos = append(photos, photo.Path)
	}

	itemResponse.ID = item.ID
	itemResponse.Code = item.Code
	itemResponse.Title = item.Title
	itemResponse.Description = item.Description
	itemResponse.Price = item.Price
	itemResponse.Stock = item.Stock
	itemResponse.ItemType = item.ItemType
	itemResponse.Leader = item.Leader
	itemResponse.LeaderLevel = item.LeaderLevel
	itemResponse.Status = item.Status
	itemResponse.Photos = photos
	itemResponse.CreatedAt = item.CreatedAt
	itemResponse.UpdatedAt = item.UpdatedAt
}

type ErrorResponse struct {
	Message string `json:"message"`
}
