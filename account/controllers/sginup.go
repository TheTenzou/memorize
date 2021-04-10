package controllers

import (
	"log"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signinRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (c *Controller) Signup(context *gin.Context) {
	var request signinRequest

	if ok := bindData(context, &request); !ok {
		return
	}

	user := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}

	err := c.UserService.Signup(context, user)

	if err != nil {
		log.Printf("Faild to sign up user: %v\n", err.Error())
		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := c.TokenService.NewPairFromUser(context, user, "")

	if err != nil {
		log.Printf("Failded to create tokens for user: %v\n", err.Error())

		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})

		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
