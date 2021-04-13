package middleware

import (
	"memorize/models"
	"memorize/models/apperrors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// extracts a user from the Authorization header
// It sets the user to the context if the user exists
func AuthUser(s models.TokenService) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		header := authHeader{}

		if err := ginContext.ShouldBindHeader(&header); err != nil {
			handleError(ginContext, err)
			ginContext.Abort()
			return
		}

		idTokenHeader := strings.Split(header.Token, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			ginContext.JSON(err.Status(), gin.H{
				"error": err,
			})
			ginContext.Abort()
			return
		}

		user, err := s.ValidateAccessToken(idTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("Provided token is invalid")
			ginContext.JSON(err.Status(), gin.H{
				"error": err,
			})
			ginContext.Abort()
			return
		}

		ginContext.Set("user", user)

		ginContext.Next()
	}
}

func handleError(ginContext *gin.Context, err error) {
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

		err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

		ginContext.JSON(err.Status(), gin.H{
			"error":       err,
			"invalidArgs": invalidArgs,
		})
		ginContext.Abort()
		return
	}

	error := apperrors.NewInternal()
	ginContext.JSON(error.Status(), gin.H{
		"error": error,
	})
}
