package controller

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

func (c *controller) Signup(ctx *gin.Context) {
	var request signinRequest

	if ok := bindData(ctx, &request); !ok {
		return
	}

	user := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}

	requestContext := ctx.Request.Context()
	err := c.UserService.Signup(requestContext, user)

	if err != nil {
		log.Printf("Faild to sign up user: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := c.TokenService.NewPairFromUser(requestContext, user, "")

	if err != nil {
		log.Printf("Failded to create tokens for user: %v\n", err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
