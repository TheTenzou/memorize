package controllers

import (
	"log"
	"memorize/models"
	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
)

type signinRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (c *Controller) Signup(ctx *gin.Context) {
	var request signinRequest

	if ok := bindData(ctx, &request); !ok {
		return
	}

	user := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}

	err := c.UserService.Signup(ctx, user)

	if err != nil {
		log.Printf("Faild to sign up user: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
}
