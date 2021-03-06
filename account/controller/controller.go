package controller

import (
	"net/http"
	"time"

	"memorize/controller/middleware"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
)

type controller struct {
	UserService  models.UserService
	TokenService models.TokenService
}

// hold services that will eventually be injected into this handler layer on handler initialization
type Config struct {
	Router          *gin.Engine
	UserService     models.UserService
	TokenService    models.TokenService
	BaseURL         string
	TImeoutDuration time.Duration
}

// initializes the handler with required injected services along with http routes
func NewController(config *Config) {

	ctrl := &controller{
		UserService:  config.UserService,
		TokenService: config.TokenService,
	}

	group := config.Router.Group(config.BaseURL)

	if gin.Mode() != gin.TestMode {
		group.Use(middleware.Timeout(config.TImeoutDuration, apperrors.NewServiceUnavailable()))
		group.GET("/me", middleware.AuthUser(ctrl.TokenService), ctrl.Me)
		group.POST("/signout", middleware.AuthUser(ctrl.TokenService), ctrl.Signout)
		group.PUT("/details", middleware.AuthUser(ctrl.TokenService), ctrl.Details)
	} else {
		group.GET("/me", ctrl.Me)
		group.POST("/signout", ctrl.Signout)
		group.PUT("/details", ctrl.Details)
	}

	group.POST("/signup", ctrl.Signup)
	group.POST("/signin", ctrl.Signin)
	group.POST("/tokens", ctrl.Tokens)
	group.POST("/image", ctrl.Image)
	group.DELETE("/image", ctrl.DeleteImage)
}

func (c *controller) Image(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

func (c *controller) DeleteImage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"hello": "it's delete image",
	})
}
