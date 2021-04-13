package controller

import (
	"log"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokensRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (c *controller) Tokens(ginContext *gin.Context) {
	var request tokensRequest

	if ok := bindData(ginContext, &request); !ok {
		return
	}

	ctx := ginContext.Request.Context()

	refreshToken, err := c.TokenService.ValidateRefreshToken(request.RefreshToken)

	if err != nil {
		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	user, err := c.UserService.GetUser(ctx, refreshToken.UserID)

	if err != nil {
		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := c.TokenService.NewPairFromUser(ctx, user, refreshToken.ID.String())

	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", user, err.Error())

		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
