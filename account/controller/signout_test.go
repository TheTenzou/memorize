package controller

import (
	"encoding/json"
	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignout(test *testing.T) {
	gin.SetMode(gin.TestMode)

	test.Run("Success", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		ctxUser := &models.User{
			UID:   uid,
			Login: "alice",
		}

		recorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", ctxUser)
		})

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.On("Signout", mock.AnythingOfType("*context.emptyCtx"), ctxUser.UID).Return(nil)

		NewController(&Config{
			Router:       router,
			TokenService: mockTokenService,
		})

		request, _ := http.NewRequest(http.MethodPost, "/signout", nil)
		router.ServeHTTP(recorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"message": "user signed out successfully!",
		})

		assert.Equal(test, http.StatusOK, recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())
	})

	test.Run("Signout Error", func(test *testing.T) {
		uid, _ := uuid.NewRandom()

		ctxUser := &models.User{
			UID:   uid,
			Login: "alice",
		}

		recorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", ctxUser)
		})

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.
			On("Signout", mock.AnythingOfType("*context.emptyCtx"), ctxUser.UID).
			Return(apperrors.NewInternal())

		NewController(&Config{
			Router:       router,
			TokenService: mockTokenService,
		})

		request, _ := http.NewRequest(http.MethodPost, "/signout", nil)
		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusInternalServerError, recorder.Code)
	})
}
