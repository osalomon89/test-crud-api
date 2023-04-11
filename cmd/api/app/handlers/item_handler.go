package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/osalomon89/test-crud-api/cmd/api/app/handlers/presenter"
	"github.com/osalomon89/test-crud-api/internal/market"
	"github.com/osalomon89/test-crud-api/internal/market/item"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type ItemHandler interface {
	CreateItem(c *gin.Context)
	GetItemByID(c *gin.Context)
}

type itemHandler struct {
	itemUsecase market.CRUDUseCase
}

func NewItemHandler(useCase market.CRUDUseCase) (ItemHandler, error) {
	if useCase == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	return &itemHandler{
		itemUsecase: useCase,
	}, nil
}

func (h *itemHandler) CreateItem(c *gin.Context) {
	ctx := marketcontext.New(c.Request)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering ItemHandler CreateItem()")

	var itemBody itemRequest
	if err := c.BindJSON(&itemBody); err != nil {
		c.JSON(http.StatusBadRequest, presenter.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	newItem, err := h.itemUsecase.Create(ctx, itemBody.toItemDomain())
	if err != nil {
		logger.Error(h, nil, err, "error creating item")
		var errorMsg string
		var httpStatus int

		itemError := new(item.ItemError)
		if ok := errors.As(err, itemError); ok {
			httpStatus = http.StatusBadRequest
			errorMsg = itemError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		c.JSON(httpStatus, presenter.ApiError{
			StatusCode: httpStatus,
			Message:    errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, presenter.ItemResponse{
		Response: presenter.Response{
			Status:  http.StatusCreated,
			Message: "success",
		},
		Data: presenter.Item(newItem),
	})
}

func (h *itemHandler) GetItemByID(c *gin.Context) {
	ctx := marketcontext.New(c.Request)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering ItemHandler. GetItemByID()")

	id, err := strconv.ParseUint(web.Params(c.Request)["id"], 10, 32)
	if err != nil || id <= 0 {
		logger.Error(h, nil, err, "error validating request param")

		c.JSON(http.StatusBadRequest, presenter.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	result, err := h.itemUsecase.Read(ctx, uint(id))
	if err != nil {
		logger.Error(h, nil, err, "error getting item by ID")
		var errorMsg string
		var httpStatus int

		itemError := new(item.ResourceNotFoundError)
		if ok := errors.As(err, itemError); ok {
			httpStatus = http.StatusNotFound
			errorMsg = itemError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		c.JSON(httpStatus, presenter.ApiError{
			StatusCode: httpStatus,
			Message:    errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, presenter.ItemResponse{
		Response: presenter.Response{
			Status:  http.StatusCreated,
			Message: "success",
		},
		Data: presenter.Item(result),
	})
}
