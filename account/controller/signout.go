package controller

import (
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *controller) Signout(ginContext *gin.Context) {
	user := ginContext.MustGet("user")

	ctx := ginContext.Request.Context()
	if err := c.TokenService.Signout(ctx, user.(*models.User).UID); err != nil {
		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully!",
	})
}
