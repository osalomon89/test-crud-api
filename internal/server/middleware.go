package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/osalomon89/test-crud-api/internal/server/httpcontext"
)

func newServerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New().String()
		c.Request.Header.Add(httpcontext.ReqIdHeaderName, uuid)
		c.Next()
	}
}
