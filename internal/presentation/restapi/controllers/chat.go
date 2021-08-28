package controllers

import (
	"encoding/json"
	"fmt"
	"time"

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

type message struct {
	From   string `json:"from"`
	Body   string `json:"body"`
	SentAt int64  `json:"sent_at"`
}

func (ctrl *ChatController) HandleWsMessage(s *melody.Session, msg []byte) {
	logger.Log(logger.DEBUG, fmt.Sprintf("new message %v : %s", s.Keys, string(msg)))

	sender, ok := s.Get(SENDER)
	if !ok {
		logger.Log(logger.ERROR, fmt.Sprintf("session %v doesn't have sender id", s))
		sender = interface{}("")
	}
	message := message{
		From:   sender.(string),
		Body:   string(msg),
		SentAt: time.Now().UTC().Unix(),
	}

	jmsg, err := json.Marshal(message)
	if err != nil {
		logger.Log(logger.ERROR, fmt.Sprintf("failed to marshal %v message\nerr: %v", message, err))
		return
	}

	ctrl.ws.BroadcastFilter(jmsg, func(q *melody.Session) bool {
		// do not send msg to our current session
		if s == q {
			return false
		}
		// get our id to send message to our other online devices
		ourId, ok := s.Get(SENDER)
		if !ok {
			logger.Log(logger.WARN, fmt.Sprintf("session %v doesn't have sender id", s))
			return false
		}
		// get our channel destination user id which we want to send a message
		ourDestinationUserId, ok := s.Get(RECEIVER)
		if !ok {
			logger.Log(logger.WARN, fmt.Sprintf("session %v doesn't have receiver id", s))
			return false
		}
		// get other user id from their session to compare with our destination user
		otherUserId, ok := q.Get(SENDER)
		if !ok {
			logger.Log(logger.WARN, fmt.Sprintf("session %v doesn't have sender id", s))
			return false
		}

		return ourDestinationUserId == otherUserId || otherUserId == ourId
	})
}
