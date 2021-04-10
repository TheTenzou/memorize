package controllers

import (
	"net/http"

	"memorize/models"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	UserService  models.UserService
	TokenService models.TokenService
}

type Config struct {
	Router       *gin.Engine
	UserService  models.UserService
	TokenService models.TokenService
}

func NewController(config *Config) {

	controller := &Controller{
		UserService:  config.UserService,
		TokenService: config.TokenService,
	}

	g := config.Router.Group("/api/account")

	g.GET("/me", controller.Me)
	g.POST("/signup", controller.Signup)
	g.POST("/signin", controller.Signin)
	g.POST("/sigout", controller.Signout)
	g.POST("/tokens", controller.Tokens)
	g.POST("/image", controller.Image)
	g.DELETE("/image", controller.DeleteImage)
	g.PUT("/details", controller.Details)
}

func (c *Controller) Signin(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

func (c *Controller) Signout(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

func (c *Controller) Tokens(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

func (c *Controller) Image(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

func (c *Controller) DeleteImage(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"hello": "it's delete image",
	})
}

func (c *Controller) Details(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"cello": "it's details",
	})
}
