package restapi

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) registerRoutes(r *gin.Engine) {

	r.POST("/register", s.controllers.Auth.Register)
	r.POST("/login", s.controllers.Auth.Login)

	users := r.Group("/users", s.middlewares.Auth.JWT)
	users.GET("/", s.controllers.User.ListUsers)

	chats := r.Group("/chats", s.middlewares.Auth.JWT)
	chats.GET("/:id", s.controllers.Chat.StartChat)
}
