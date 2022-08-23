package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/osalomon89/test-crud-api/internal/application/ports"
	"github.com/osalomon89/test-crud-api/internal/domain"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type ItemHandler interface {
	CreateItem(res http.ResponseWriter, req *http.Request) error
	GetItemByID(res http.ResponseWriter, req *http.Request) error
}

type itemHandler struct {
	itemService ports.ItemService
}

func NewItemHandler(itemService ports.ItemService) ItemHandler {
	return &itemHandler{
		itemService: itemService,
	}
}

func (h *itemHandler) CreateItem(res http.ResponseWriter, req *http.Request) error {
	ctx := marketcontext.New(req)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering ItemHandler CreateItem()")

	var itemBody domain.CreateItemBody
	if err := json.NewDecoder(req.Body).Decode(&itemBody); err != nil {
		logger.Error(h, nil, err, "error validating request body")

		message := Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}

		return web.EncodeJSON(res, message, http.StatusBadRequest)
	}

	item, err := h.itemService.CreateItem(ctx, itemBody)
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

		message := Response{
			Status:  httpStatus,
			Message: errorMsg,
			Data:    nil,
		}

		return web.EncodeJSON(res, message, http.StatusBadRequest)
	}

	return web.EncodeJSON(res, Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    createItemResponse(item),
	}, http.StatusCreated)
}

func (h *itemHandler) GetItemByID(res http.ResponseWriter, req *http.Request) error {
	ctx := marketcontext.New(req)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering ItemHandler. GetItemByID()")

	id, err := strconv.ParseUint(web.Params(req)["id"], 10, 32)
	if err != nil || id <= 0 {
		logger.Error(h, nil, err, "error validating request param")

		return web.EncodeJSON(res, Response{
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

		return web.EncodeJSON(res, Response{
			Status:  httpStatus,
			Message: errorMsg,
			Data:    nil,
		}, httpStatus)
	}

	return web.EncodeJSON(res, Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    createItemResponse(item),
	}, http.StatusOK)
}
