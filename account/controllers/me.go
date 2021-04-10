package controllers

import (
	"log"
	"net/http"

	"memorize/models"

	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
)

func (c *controller) Me(context *gin.Context) {
	user, exists := context.Get("user")

	if !exists {
		log.Printf("Unable to extract user from request context unknown reason: %v\n", context)
		err := apperrors.NewInternal()
		context.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*models.User).UID

	u, err := c.UserService.Get(context, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		context.JSON(e.Status(), gin.H{
			"error": e,
		})

		return
	}

	context.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
