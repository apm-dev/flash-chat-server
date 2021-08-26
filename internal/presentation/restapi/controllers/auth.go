package controllers

import (
	"errors"
	"net/http"

	"github.com/apm-dev/flash-chat/internal/domain"
	"github.com/apm-dev/flash-chat/internal/domain/authing"
	request "github.com/apm-dev/flash-chat/internal/presentation/restapi/requests"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi/responses"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc authing.Service
}

func NewAuthController(svc authing.Service) *AuthController {
	return &AuthController{svc}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var data request.Register
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.Make(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	token, err := ctrl.svc.Register(data.Name, data.Username, data.Password)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, domain.ErrInternalServer) {
			code = http.StatusInternalServerError
		}
		c.JSON(code, responses.Make(code, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, responses.Make(
		http.StatusOK, "welcome", gin.H{"token": token},
	))
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var data request.Login
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.Make(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	token, err := ctrl.svc.Login(data.Username, data.Password)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, domain.ErrInternalServer) {
			code = http.StatusInternalServerError
		}
		c.JSON(code, responses.Make(code, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, responses.Make(http.StatusOK, "", gin.H{"token": token}))
}
