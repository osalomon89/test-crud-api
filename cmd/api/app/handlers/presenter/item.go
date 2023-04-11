package presenter

import (
	"time"

	"github.com/osalomon89/test-crud-api/internal/market/item"
)

type jsonItem struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	ItemType    string    `json:"itemType"`
	Leader      bool      `json:"leader"`
	LeaderLevel string    `json:"leaderLevel"`
	Status      string    `json:"status"`
	Photos      []string  `json:"photos"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func Item(i *item.Item) *jsonItem {
	var itemResponse jsonItem

	itemResponse.ID = i.ID
	itemResponse.Code = i.Code
	itemResponse.Title = i.Title
	itemResponse.Description = i.Description
	itemResponse.Price = i.Price
	itemResponse.Stock = i.Stock
	itemResponse.ItemType = i.ItemType
	itemResponse.Leader = i.Leader
	itemResponse.LeaderLevel = i.LeaderLevel
	itemResponse.Status = i.Status
	itemResponse.Photos = i.Photos
	itemResponse.CreatedAt = i.CreatedAt
	itemResponse.UpdatedAt = i.UpdatedAt

	return &itemResponse
}
