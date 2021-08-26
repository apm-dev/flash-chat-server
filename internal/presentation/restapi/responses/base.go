package responses

import (
	"strconv"
)

type Base struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

func Make(code int, message string, content interface{}) Base {
	return Base{
		Code:    strconv.Itoa(code),
		Message: message,
		Content: content,
	}
}
