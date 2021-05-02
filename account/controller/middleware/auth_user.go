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
	return func(ctx *gin.Context) {
		header := authHeader{}

		if err := ctx.ShouldBindHeader(&header); err != nil {
			handleError(ctx, err)
			ctx.Abort()
			return
		}

		idTokenHeader := strings.Split(header.Token, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		user, err := s.ValidateAccessToken(idTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("Provided token is invalid")
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}

func handleError(ctx *gin.Context, err error) {
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

		ctx.JSON(err.Status(), gin.H{
			"error":       err,
			"invalidArgs": invalidArgs,
		})
		ctx.Abort()
		return
	}

	error := apperrors.NewInternal()
	ctx.JSON(error.Status(), gin.H{
		"error": error,
	})
}
