package controllers

import (
	"net/http"

	"memorize/models"

	"github.com/gin-gonic/gin"
)

type controller struct {
	UserService  models.UserService
	TokenService models.TokenService
}

type Config struct {
	Router       *gin.Engine
	UserService  models.UserService
	TokenService models.TokenService
	BaseURL      string
}

func NewController(config *Config) {

	ctrl := &controller{
		UserService:  config.UserService,
		TokenService: config.TokenService,
	}

	g := config.Router.Group(config.BaseURL)

	g.GET("/me", ctrl.Me)
	g.POST("/signup", ctrl.Signup)
	g.POST("/signin", ctrl.Signin)
	g.POST("/sigout", ctrl.Signout)
	g.POST("/tokens", ctrl.Tokens)
	g.POST("/image", ctrl.Image)
	g.DELETE("/image", ctrl.DeleteImage)
	g.PUT("/details", ctrl.Details)
}

func (c *controller) Signin(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

func (c *controller) Signout(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

func (c *controller) Tokens(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

func (c *controller) Image(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

func (c *controller) DeleteImage(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's delete image",
	})
}

func (c *controller) Details(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"cello": "it's details",
	})
}
