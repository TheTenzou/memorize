package controllers

import (
	"net/http"

	"memorize/models"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	UserService models.UserService
}

type Config struct {
	Router      *gin.Engine
	UserService models.UserService
}

func NewController(c *Config) {

	h := &Controller{
		UserService: c.UserService,
	}

	g := c.Router.Group("/api/account")

	g.GET("/me", h.Me)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/sigout", h.Signout)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	g.PUT("/details", h.Details)
}

func (h *Controller) Signin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

func (h *Controller) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

func (h *Controller) Tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

func (h *Controller) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

func (h *Controller) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's delete image",
	})
}

func (h *Controller) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's details",
	})
}
