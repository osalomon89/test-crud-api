package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/osalomon89/test-crud-api/internal/core/domain"
	"github.com/osalomon89/test-crud-api/internal/core/ports"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server/handler/dto"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type ItemHandler interface {
	CreateItem(res http.ResponseWriter, req *http.Request) error
	GetItemByID(res http.ResponseWriter, req *http.Request) error
}

type itemHandler struct {
	itemService ports.ItemService
}

func NewItemHandler(itemService ports.ItemService) (ItemHandler, error) {
	if itemService == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	return &itemHandler{
		itemService: itemService,
	}, nil
}

func (h *itemHandler) CreateItem(res http.ResponseWriter, req *http.Request) error {
	ctx := marketcontext.New(req)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering ItemHandler CreateItem()")

	var itemBody dto.ItemBody
	if err := json.NewDecoder(req.Body).Decode(&itemBody); err != nil {
		logger.Error(h, nil, err, "error validating request body")

		message := dto.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}

		return web.EncodeJSON(res, message, http.StatusBadRequest)
	}

	item, err := h.itemService.CreateItem(ctx, itemBody.ToItemDomain())
	if err != nil {
		logger.Error(h, nil, err, "error creating item")
		var errorMsg string
		var httpStatus int

		itemError := new(domain.ItemError)
		if ok := errors.As(err, itemError); ok {
			httpStatus = http.StatusBadRequest
			errorMsg = itemError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		message := dto.Response{
			Status:  httpStatus,
			Message: errorMsg,
			Data:    nil,
		}

		return web.EncodeJSON(res, message, httpStatus)
	}

	return web.EncodeJSON(res, dto.Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    dto.CreateItemResponse(item),
	}, http.StatusCreated)
}

func (h *itemHandler) GetItemByID(res http.ResponseWriter, req *http.Request) error {
	ctx := marketcontext.New(req)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering ItemHandler. GetItemByID()")

	id, err := strconv.ParseUint(web.Params(req)["id"], 10, 32)
	if err != nil || id <= 0 {
		logger.Error(h, nil, err, "error validating request param")

		return web.EncodeJSON(res, dto.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}, http.StatusBadRequest)
	}

	item, err := h.itemService.GetItemByID(ctx, uint(id))
	if err != nil {
		logger.Error(h, nil, err, "error getting item by ID")
		var errorMsg string
		var httpStatus int

		itemError := new(domain.ResourceNotFoundError)
		if ok := errors.As(err, itemError); ok {
			httpStatus = http.StatusNotFound
			errorMsg = itemError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		return web.EncodeJSON(res, dto.Response{
			Status:  httpStatus,
			Message: errorMsg,
			Data:    nil,
		}, httpStatus)
	}

	return web.EncodeJSON(res, dto.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    dto.CreateItemResponse(item),
	}, http.StatusOK)
}
