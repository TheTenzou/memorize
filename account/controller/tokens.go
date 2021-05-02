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

func (c *controller) Tokens(ctx *gin.Context) {
	var request tokensRequest

	if ok := bindData(ctx, &request); !ok {
		return
	}

	requestCtx := ctx.Request.Context()

	refreshToken, err := c.TokenService.ValidateRefreshToken(request.RefreshToken)

	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	user, err := c.UserService.GetUser(requestCtx, refreshToken.UserID)

	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := c.TokenService.NewPairFromUser(requestCtx, user, refreshToken.ID.String())

	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", user, err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
