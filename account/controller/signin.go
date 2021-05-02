package controller

import (
	"log"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// signinRequset is not exported
type signinRequset struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signin used to authenticate extant user
func (c *controller) Signin(ctx *gin.Context) {
	var request signinRequset

	if ok := bindData(ctx, &request); !ok {
		return
	}

	user := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}

	requestCtx := ctx.Request.Context()
	signupedUser, err := c.UserService.Signin(requestCtx, user)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := c.TokenService.NewPairFromUser(requestCtx, signupedUser, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
