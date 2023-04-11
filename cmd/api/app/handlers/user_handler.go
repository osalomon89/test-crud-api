package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osalomon89/test-crud-api/cmd/api/app/handlers/presenter"
	"github.com/osalomon89/test-crud-api/internal/market"
	"github.com/osalomon89/test-crud-api/internal/market/user"
	marketcontext "github.com/osalomon89/test-crud-api/pkg/context"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	MarkItemAsFavorite(c *gin.Context)
}

type userHandler struct {
	userUsecase market.UserUseCase
}

func NewUserHandler(useCase market.UserUseCase) (UserHandler, error) {
	if useCase == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	return &userHandler{
		userUsecase: useCase,
	}, nil
}

type registerRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,min=8"`
}

func (h *userHandler) Register(c *gin.Context) {
	ctx := marketcontext.New(c.Request)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering UserHandler SaveUser()")

	var payload registerRequest
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, presenter.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		c.JSON(http.StatusBadRequest, presenter.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Passwords does not match",
		})
		return
	}

	err := h.userUsecase.Register(ctx, payload.Email, payload.Password)
	if err != nil {
		logger.Error(h, nil, err, "error creating user")
		var errorMsg string
		var httpStatus int

		userError := new(user.UserError)
		if ok := errors.As(err, userError); ok {
			httpStatus = http.StatusBadRequest
			errorMsg = userError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		c.JSON(httpStatus, presenter.UserResponse{
			Response: presenter.Response{
				Status:  httpStatus,
				Message: errorMsg,
			},
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusCreated, presenter.UserResponse{
		Response: presenter.Response{
			Status:  http.StatusCreated,
			Message: "success",
		},
		Data: presenter.User(payload.Email),
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *userHandler) Login(c *gin.Context) {
	ctx := marketcontext.New(c.Request)
	logger := marketcontext.Logger(ctx)
	logger.Debug(h, nil, "Entering UserHandler. Login()")

	var payload loginRequest
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, presenter.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	result, err := h.userUsecase.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		logger.Error(h, nil, err, "error getting item by ID")
		var errorMsg string
		var httpStatus int

		userError := new(user.ResourceNotFoundError)
		if ok := errors.As(err, userError); ok {
			httpStatus = http.StatusNotFound
			errorMsg = userError.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			errorMsg = err.Error()
		}

		c.JSON(httpStatus, presenter.UserResponse{
			Response: presenter.Response{
				Status:  httpStatus,
				Message: errorMsg,
			},
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, presenter.TokenResponse{
		Response: presenter.Response{
			Status:  http.StatusOK,
			Message: "success",
		},
		AccessToken: result,
	})
}

func (h *userHandler) MarkItemAsFavorite(c *gin.Context) {
	userID := c.GetHeader("userID")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s: pong", userID)})
}
