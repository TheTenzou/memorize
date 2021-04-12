package inject

import (
	"memorize/controller"
	"os"

	"github.com/gin-gonic/gin"
)

func InitRouter(services *Services) *gin.Engine {

	router := gin.Default()

	baseUrl := os.Getenv("ACCOUNT_API_URL")
	controller.NewController(&controller.Config{
		Router:       router,
		UserService:  services.UserService,
		TokenService: services.TokenService,
		BaseURL:      baseUrl,
	})

	return router
}
