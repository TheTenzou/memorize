package controller

import (
	"log"
	"net/http"

	"memorize/models"

	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
)

func (c *controller) Me(ctx *gin.Context) {
	user, exists := ctx.Get("user")

	if !exists {
		log.Printf("Unable to extract user from request context unknown reason: %v\n", ctx)
		err := apperrors.NewInternal()
		ctx.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*models.User).UID

	requestContext := ctx.Request.Context()
	userFromDB, err := c.UserService.GetUser(requestContext, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		ctx.JSON(e.Status(), gin.H{
			"error": e,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": userFromDB,
	})
}
