package inject

import (
	"fmt"
	"memorize/controller"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Init router, inject servisec into controller
func InitRouter(services *Services) (*gin.Engine, error) {

	router := gin.Default()

	baseUrl := os.Getenv("ACCOUNT_API_URL")

	controllerTimeoutStr := os.Getenv("CONTROLLER_TIMEOUT")
	controllerTimeout, err := strconv.ParseInt(controllerTimeoutStr, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	controller.NewController(&controller.Config{
		Router:          router,
		UserService:     services.UserService,
		TokenService:    services.TokenService,
		BaseURL:         baseUrl,
		TImeoutDuration: time.Duration(time.Duration(controllerTimeout) * time.Second),
	})

	return router, nil
}
