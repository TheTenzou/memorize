package controller

import (
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *controller) Signout(ctx *gin.Context) {
	user := ctx.MustGet("user")

	requestCtx := ctx.Request.Context()
	if err := c.TokenService.Signout(requestCtx, user.(*models.User).UID); err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully!",
	})
}
