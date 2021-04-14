package controller

import (
	"log"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type detailsReq struct {
	Name    string `json:"name" binding:"omitempty,max=50"`
	Email   string `json:"email" binding:"omitempty,email"`
	Website string `json:"website" binding:"omitempty,url"`
}

func (c *controller) Details(ginContext *gin.Context) {
	authUser := ginContext.MustGet("user").(*models.User)

	var request detailsReq

	if ok := bindData(ginContext, &request); !ok {
		return
	}

	user := &models.User{
		UID:     authUser.UID,
		Name:    request.Name,
		Email:   request.Email,
		Website: request.Website,
	}

	ctx := ginContext.Request.Context()
	err := c.UserService.UpdateDetails(ctx, user)

	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		ginContext.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
