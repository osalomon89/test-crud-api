package httpmodel

import "github.com/osalomon89/test-crud-api/internal/model"

type CreateItemBody struct {
	Code        string   `json:"code" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       int      `json:"price" binding:"required"`
	Stock       int      `json:"stock" binding:"required"`
	ItemType    string   `json:"itemType" binding:"required"`
	Leader      bool     `json:"leader"`
	LeaderLevel string   `json:"leaderLevel"`
	Status      string   `json:"status"`
	Photos      []string `json:"photos" binding:"required"`
}

func (p CreateItemBody) GetStatus() string {
	if p.Stock > 0 {
		return model.StatusActive
	}

	return model.StatusInactive
}
