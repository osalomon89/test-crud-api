package server

import (
	"time"

	"github.com/osalomon89/test-crud-api/internal/domain"
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

func createItemResponse(item *domain.Item) *ItemResponse {
	var photos []string
	for _, photo := range item.Photos {
		photos = append(photos, photo.Path)
	}

	return &ItemResponse{
		ID:          item.ID,
		Code:        item.Code,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		ItemType:    item.ItemType,
		Leader:      item.Leader,
		LeaderLevel: item.LeaderLevel,
		Status:      item.Status,
		Photos:      photos,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
