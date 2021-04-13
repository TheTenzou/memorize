package middleware

import (
	"fmt"
	"memorize/mocks"
	"memorize/models"
	"memorize/models/apperrors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthUser(test *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(mocks.MockTokenService)

	uid, _ := uuid.NewRandom()
	user := &models.User{
		UID:   uid,
		Login: "alice",
	}

	validTokenHeader := "validTokenString"
	invalidTokenHeader := "invalidTokenString"
	invalidTokenErr := apperrors.NewAuthorization("Unable to verify user from idToken")

	mockTokenService.On("ValidateAccessToken", validTokenHeader).Return(user, nil)
	mockTokenService.On("ValidateAccessToken", invalidTokenHeader).Return(nil, invalidTokenErr)

	test.Run("Adds a user to context", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		_, router := gin.CreateTestContext(recorder)

		var contextUser *models.User

		router.GET("/me", AuthUser(mockTokenService), func(c *gin.Context) {
			contextKeyVal, _ := c.Get("user")
			contextUser = contextKeyVal.(*models.User)
		})

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validTokenHeader))
		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusOK, recorder.Code)
		assert.Equal(test, user, contextUser)

		mockTokenService.AssertCalled(test, "ValidateAccessToken", validTokenHeader)
	})

	test.Run("Invalid Token", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		_, router := gin.CreateTestContext(recorder)

		router.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidTokenHeader))
		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusUnauthorized, recorder.Code)
		mockTokenService.AssertCalled(test, "ValidateAccessToken", invalidTokenHeader)
	})

	test.Run("Missing Authorization Header", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		_, router := gin.CreateTestContext(recorder)

		router.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusUnauthorized, recorder.Code)
		mockTokenService.AssertNotCalled(test, "ValidateAccessToken")
	})
}
