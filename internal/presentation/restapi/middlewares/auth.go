package middlewares

import (
	"net/http"
	"strings"

	"github.com/apm-dev/flash-chat/internal/domain/authing"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi/responses"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authing authing.Service
}

func NewAuthMiddleware(auth authing.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authing: auth,
	}
}

func (m *AuthMiddleware) JWT(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.Make(
			http.StatusUnauthorized,
			"Authorization header is required",
			nil,
		))
		return
	}

	token = strings.TrimSpace(strings.Replace(token, "Bearer", "", 1))
	user, err := m.authing.Authorize(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.Make(
			http.StatusUnauthorized,
			err.Error(),
			nil,
		))
		return
	}

	c.Set("user", user)
}
