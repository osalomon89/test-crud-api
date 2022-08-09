package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	itemerrors "github.com/osalomon89/test-crud-api/internal/errors"
	httpmodel "github.com/osalomon89/test-crud-api/internal/model/http"
	"github.com/osalomon89/test-crud-api/internal/server/httpcontext"
	"github.com/osalomon89/test-crud-api/internal/service"
)

type ItemHTTPHandler interface {
	CreateItem(ctx *gin.Context)
	GetItemByID(ctx *gin.Context)
}

// ItemHTTPHandler represents article http handler
type itemHTTPHandler struct {
	itemService service.ItemService
}

// NewItemHTTPHandler return new instances of item http handler
func NewItemHTTPHandler(itemService service.ItemService) ItemHTTPHandler {
	return &itemHTTPHandler{
		itemService: itemService,
	}
}

func (h *itemHTTPHandler) CreateItem(c *gin.Context) {
	ctx := httpcontext.BackgroundFromContext(c)
	logger := httpcontext.GetLogger(ctx)

	logger.Info("Entering ItemHandler. CreateItem()")

	var itemBody httpmodel.CreateItemBody
	if err := c.ShouldBindJSON(&itemBody); err != nil {
		logger.Debug("Error validating request body. ", err)

		c.JSON(http.StatusBadRequest, &httpmodel.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	item, err := h.itemService.CreateItem(ctx, itemBody)
	if err != nil {
		var errorMsg string
		var httpStatus int

		itemError := new(itemerrors.ItemError)
		if ok := errors.As(err, itemError); ok {
			httpStatus = http.StatusBadRequest
			errorMsg = itemError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		c.JSON(httpStatus, &httpmodel.Response{
			Status:  httpStatus,
			Message: errorMsg,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, &httpmodel.Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    httpmodel.CreateItemResponse(item),
	})
}

func (h *itemHTTPHandler) GetItemByID(c *gin.Context) {
	ctx := httpcontext.BackgroundFromContext(c)
	logger := httpcontext.GetLogger(ctx)

	logger.Info("Entering ItemHandler. GetItemByID()")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id <= 0 {
		logger.Debug("Error validating request param. ", err)

		c.JSON(http.StatusBadRequest, &httpmodel.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	item, err := h.itemService.GetItemByID(ctx, uint(id))
	if err != nil {
		var errorMsg string
		var httpStatus int

		itemError := new(itemerrors.ResourceNotFoundError)
		if ok := errors.As(err, itemError); ok {
			httpStatus = http.StatusNotFound
			errorMsg = itemError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		c.JSON(httpStatus, &httpmodel.Response{
			Status:  httpStatus,
			Message: errorMsg,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, &httpmodel.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    httpmodel.CreateItemResponse(item),
	})
}
