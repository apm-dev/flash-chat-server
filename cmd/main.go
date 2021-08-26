package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/apm-dev/flash-chat/internal/data"
	"github.com/apm-dev/flash-chat/internal/data/storage/memory"
	"github.com/apm-dev/flash-chat/internal/domain/authing"
	"github.com/apm-dev/flash-chat/internal/domain/listing"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi/controllers"
	"github.com/apm-dev/flash-chat/internal/presentation/restapi/middlewares"
	"github.com/apm-dev/flash-chat/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger.SetLogger(logger.NewLogcat(logger.DEBUG))

	// configs, env variables
	err := godotenv.Load()
	if err != nil {
		logger.Log(logger.WARN, "couldn't find .env file")
	}

	// DataBases, MQs, APIs client
	userDS := memory.NewUserDS()

	// Repositories
	userRepo := data.NewUserRepo(userDS)

	// Pseudo services
	jwtMng := authing.NewJWTManager(getJwtKeyAndExp())

	// Services
	authSvc := authing.NewService(userRepo, jwtMng)
	listingSvc := listing.NewService(userRepo)

	// Start
	httpServeAddr := getEnvOrDefault("HTTP_HOST", ":8080")
	// Handlers(controllers,entry points)
	ctrls := restapi.Controllers{
		Auth: controllers.NewAuthController(authSvc),
		Chat: controllers.NewChatController(),
		User: controllers.NewUserController(listingSvc),
	}

	midlwrs := restapi.Middlewares{
		Auth: middlewares.NewAuthMiddleware(authSvc),
	}
	// Servers
	restSrv := restapi.NewServer(ctrls, midlwrs)

	// Start
	restSrv.Start(httpServeAddr)

	// Listen to OS interrupt signal to stop app
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Stop
	restSrv.Stop()
}

func getJwtKeyAndExp() (string, time.Duration) {
	jwtKey := getEnvOrDefault("JWT_KEY", "y0ur-jwt-Encrypt!on-k3y")
	jwtExpStr := getEnvOrDefault("JWT_EXP", "120")

	jwtExpInt, err := strconv.ParseInt(jwtExpStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return jwtKey, time.Minute * time.Duration(jwtExpInt)
}

func getEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}
