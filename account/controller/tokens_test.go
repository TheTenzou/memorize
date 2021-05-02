package controller

import (
	"bytes"
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

func TestTokens(test *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(mocks.MockTokenService)
	mockUserService := new(mocks.MockUserService)

	router := gin.Default()

	NewController(&Config{
		Router:       router,
		TokenService: mockTokenService,
		UserService:  mockUserService,
	})

	test.Run("Invalid request", func(test *testing.T) {
		recorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(gin.H{
			"notRefreshToken": "this key is not valid for this handler!",
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		assert.Equal(test, http.StatusBadRequest, recorder.Code)
		mockTokenService.AssertNotCalled(test, "ValidateRefreshToken")
		mockUserService.AssertNotCalled(test, "GetUser")
		mockTokenService.AssertNotCalled(test, "NewPairFromUser")
	})

	test.Run("Invalid token", func(test *testing.T) {
		invalidTokenString := "invalid"
		mockErrorMessage := "authProbs"
		mockError := apperrors.NewAuthorization(mockErrorMessage)

		mockTokenService.
			On("ValidateRefreshToken", invalidTokenString).
			Return(nil, mockError)

		recorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(gin.H{
			"refreshToken": invalidTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(test, mockError.Status(), recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())
		mockTokenService.AssertCalled(test, "ValidateRefreshToken", invalidTokenString)
		mockUserService.AssertNotCalled(test, "GetUser")
		mockTokenService.AssertNotCalled(test, "NewPairFromUser")
	})

	test.Run("Failure to create new token pair", func(test *testing.T) {
		validTokenString := "valid"
		mockTokenID, _ := uuid.NewRandom()
		mockUserID, _ := uuid.NewRandom()

		mockRefreshTokenResponse := &models.RefreshToken{
			Token:  validTokenString,
			ID:     mockTokenID,
			UserID: mockUserID,
		}

		mockTokenService.
			On("ValidateRefreshToken", validTokenString).
			Return(mockRefreshTokenResponse, nil)

		mockUserResponse := &models.User{
			UID: mockUserID,
		}
		getArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockRefreshTokenResponse.UserID,
		}

		mockUserService.
			On("GetUser", getArguments...).
			Return(mockUserResponse, nil)

		mockError := apperrors.NewAuthorization("Invalid refresh token")
		newPairArgumnets := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUserResponse,
			mockRefreshTokenResponse.ID.String(),
		}

		mockTokenService.
			On("NewPairFromUser", newPairArgumnets...).
			Return(nil, mockError)

		recorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(gin.H{
			"refreshToken": validTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(test, mockError.Status(), recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())
		mockTokenService.AssertCalled(test, "ValidateRefreshToken", validTokenString)
		mockUserService.AssertCalled(test, "GetUser", getArguments...)
		mockTokenService.AssertCalled(test, "NewPairFromUser", newPairArgumnets...)
	})

	test.Run("Success", func(test *testing.T) {
		validTokenString := "anothervalid"
		mockTokenID, _ := uuid.NewRandom()
		mockUserID, _ := uuid.NewRandom()

		mockRefreshTokenResponse := &models.RefreshToken{
			Token:  validTokenString,
			ID:     mockTokenID,
			UserID: mockUserID,
		}

		mockTokenService.
			On("ValidateRefreshToken", validTokenString).
			Return(mockRefreshTokenResponse, nil)

		mockUserResponse := &models.User{
			UID: mockUserID,
		}
		getArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockRefreshTokenResponse.UserID,
		}

		mockUserService.
			On("GetUser", getArguments...).
			Return(mockUserResponse, nil)

		mockNewTokenID, _ := uuid.NewRandom()
		mockNewUserID, _ := uuid.NewRandom()
		mockTokenPairResponse := &models.TokenPair{
			AccessToken: models.AccessToken{Token: "aNewIDToken"},
			RefreshToken: models.RefreshToken{
				Token:  "aNewRefreshToken",
				ID:     mockNewTokenID,
				UserID: mockNewUserID,
			},
		}

		newPairArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUserResponse,
			mockRefreshTokenResponse.ID.String(),
		}

		mockTokenService.
			On("NewPairFromUser", newPairArguments...).
			Return(mockTokenPairResponse, nil)

		recorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(gin.H{
			"refreshToken": validTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"tokens": mockTokenPairResponse,
		})

		assert.Equal(test, http.StatusOK, recorder.Code)
		assert.Equal(test, responseBody, recorder.Body.Bytes())
		mockTokenService.AssertCalled(test, "ValidateRefreshToken", validTokenString)
		mockUserService.AssertCalled(test, "GetUser", getArguments...)
		mockTokenService.AssertCalled(test, "NewPairFromUser", newPairArguments...)
	})

}
