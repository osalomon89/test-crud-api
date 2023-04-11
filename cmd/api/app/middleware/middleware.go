package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/osalomon89/test-crud-api/internal/market/user/token"
)

type authMiddleware struct {
	tokenService token.Service
}

func NewAuthMiddleware(tokenService token.Service) *authMiddleware {
	return &authMiddleware{
		tokenService: tokenService,
	}
}

func (am *authMiddleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
		}

		strArr := strings.Split(authHeader, " ")
		if len(strArr) == 2 {
			tokenString = strArr[1]
		}

		tokenInfo, err := am.tokenService.Validate(c, tokenString)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Next()

			return
		}

		if tokenInfo.ID > 0 {
			userID := fmt.Sprintf("%d", tokenInfo.ID)
			c.Request.Header.Add("userID", userID)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
