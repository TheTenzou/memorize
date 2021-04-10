package controllers

import (
	"log"
	"memorize/models/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func bindData(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBind(request); err != nil {
		log.Printf("Error binding data: %v\n", err)

		handelError(c, err)

		return false
	}

	return true
}

func handelError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var invalidArgs []invalidArgument

		for _, err := range errs {
			invalidArgs = append(invalidArgs, invalidArgument{
				err.Field(),
				err.Value().(string),
				err.Tag(),
				err.Param(),
			})
		}

		err := apperrors.NewBadRequest("Invalid request parametrs. See invalidArgs")

		c.JSON(err.Status(), gin.H{
			"error":       err,
			"invalidArgs": invalidArgs,
		})

		return
	}

	fallBack := apperrors.NewInternal()

	c.JSON(fallBack.Status(), gin.H{
		"error": fallBack,
	})
}
