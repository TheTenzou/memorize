package controller

import (
	"log"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type detailsRequset struct {
	Name    string `json:"name" binding:"omitempty,max=50"`
	Email   string `json:"email" binding:"omitempty,email"`
	Website string `json:"website" binding:"omitempty,url"`
}

func (c *controller) Details(ctx *gin.Context) {
	authUser := ctx.MustGet("user").(*models.User)

	var request detailsRequset

	if ok := bindData(ctx, &request); !ok {
		return
	}

	user := &models.User{
		UID:     authUser.UID,
		Name:    request.Name,
		Email:   request.Email,
		Website: request.Website,
	}

	requestCtx := ctx.Request.Context()
	err := c.UserService.UpdateDetails(requestCtx, user)

	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
