package restapi

import (
	"fmt"
	"net/http"

	ctrl "github.com/apm-dev/flash-chat/internal/presentation/restapi/controllers"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi/middlewares"
	"github.com/apm-dev/flash-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	server      *http.Server
	quit        chan bool
	controllers Controllers
	middlewares Middlewares
}

type Controllers struct {
	Auth *ctrl.AuthController
	Chat *ctrl.ChatController
	User *ctrl.UserController
}

type Middlewares struct {
	Auth *middlewares.AuthMiddleware
}

func NewServer(ctrls Controllers, midls Middlewares) *Server {
	return &Server{
		server:      &http.Server{},
		quit:        make(chan bool, 1),
		controllers: ctrls,
		middlewares: midls,
	}
}

func (s *Server) Start(addr string) {

	fmt.Println("Starting RestServer on", addr)

	r := gin.Default()

	s.registerRoutes(r)

	s.server.Addr = addr
	s.server.Handler = r

	// start http server on different goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("RestServer stopped under request")
			} else {
				logger.Log(logger.ERROR, errors.Wrap(err, "RestServer stopped unexpectedly").Error())
				panic(err)
			}
		}
	}()

	// listen to quit channel to close the server
	go func() {
		<-s.quit
		if err := s.server.Close(); err != nil {
			logger.Log(logger.ERROR, errors.Wrap(err, "failed to stop RestServer").Error())
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("Stopping RestServer on", s.server.Addr)
	s.quit <- true
}
