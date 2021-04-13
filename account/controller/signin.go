package controller

import (
	"log"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// signinReq is not exported
type signinReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signin used to authenticate extant user
func (c *controller) Signin(ginContext *gin.Context) {
	var req signinReq

	if ok := bindData(ginContext, &req); !ok {
		return
	}

	user := &models.User{
		Login:    req.Login,
		Password: req.Password,
	}

	ctx := ginContext.Request.Context()
	signupedUser, err := c.UserService.Signin(ctx, user)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := c.TokenService.NewPairFromUser(ctx, signupedUser, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
