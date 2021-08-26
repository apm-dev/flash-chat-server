package controllers

import (
	"net/http"

	"github.com/apm-dev/flash-chat/internal/domain"
	"github.com/apm-dev/flash-chat/internal/domain/listing"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi/responses"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	listing listing.Service
}

func NewUserController(l listing.Service) *UserController {
	return &UserController{
		listing: l,
	}
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
	users, err := ctrl.listing.Users()
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, responses.Make(
			status, err.Error(), nil,
		))
		return
	}

	c.JSON(http.StatusOK, responses.Make(
		http.StatusOK, "", users,
	))
}
