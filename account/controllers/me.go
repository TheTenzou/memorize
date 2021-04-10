package controllers

import (
	"log"
	"net/http"

	"memorize/models"

	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
)

func (h *Controller) Me(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		log.Printf("Unable to extract user from request context unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*models.User).UID

	u, err := h.UserService.Get(c, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
