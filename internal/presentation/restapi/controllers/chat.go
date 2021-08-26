package controllers

import (
	"fmt"

	"github.com/apm-dev/flash-chat/internal/domain"
	"github.com/apm-dev/flash-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

const (
	SENDER   = "sender"
	RECEIVER = "receiver"
)

type ChatController struct {
	ws *melody.Melody
}

func NewChatController() *ChatController {
	ctrl := &ChatController{
		ws: melody.New(),
	}
	ctrl.ws.HandleMessage(ctrl.HandleWsMessage)
	return ctrl
}

func (ctrl *ChatController) StartChat(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	receiver := c.Param("id")
	logger.Log(logger.INFO, fmt.Sprintf(
		"user %s starts chat with %s", user.Username, receiver,
	))
	ctrl.ws.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{
		SENDER:   user.Username,
		RECEIVER: receiver,
	})
	logger.Log(logger.INFO, fmt.Sprintf(
		"user %s ends chat with %s", user.Username, receiver,
	))
}

func (ctrl *ChatController) HandleWsMessage(s *melody.Session, msg []byte) {
	logger.Log(logger.DEBUG, fmt.Sprintf("new message %v : %s", s.Keys, string(msg)))
	ctrl.ws.BroadcastFilter(msg, func(q *melody.Session) bool {
		receiver, ok := s.Get(RECEIVER)
		if !ok {
			return false
		}
		sender, ok := q.Get(SENDER)
		if !ok {
			return false
		}
		return receiver == sender
	})
}
