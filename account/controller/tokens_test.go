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

	test.Run("Invalid token", func(t *testing.T) {
		invalidTokenString := "invalid"
		mockErrorMessage := "authProbs"
		mockError := apperrors.NewAuthorization(mockErrorMessage)

		mockTokenService.
			On("ValidateRefreshToken", invalidTokenString).
			Return(nil, mockError)

		recorder := httptest.NewRecorder()

		reqBody, _ := json.Marshal(gin.H{
			"refreshToken": invalidTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())
		mockTokenService.AssertCalled(t, "ValidateRefreshToken", invalidTokenString)
		mockUserService.AssertNotCalled(t, "GetUser")
		mockTokenService.AssertNotCalled(t, "NewPairFromUser")
	})

	test.Run("Failure to create new token pair", func(t *testing.T) {
		validTokenString := "valid"
		mockTokenID, _ := uuid.NewRandom()
		mockUserID, _ := uuid.NewRandom()

		mockRefreshTokenResp := &models.RefreshToken{
			Token:  validTokenString,
			ID:     mockTokenID,
			UserID: mockUserID,
		}

		mockTokenService.
			On("ValidateRefreshToken", validTokenString).
			Return(mockRefreshTokenResp, nil)

		mockUserResp := &models.User{
			UID: mockUserID,
		}
		getArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockRefreshTokenResp.UserID,
		}

		mockUserService.
			On("GetUser", getArgs...).
			Return(mockUserResp, nil)

		mockError := apperrors.NewAuthorization("Invalid refresh token")
		newPairArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUserResp,
			mockRefreshTokenResp.ID.String(),
		}

		mockTokenService.
			On("NewPairFromUser", newPairArgs...).
			Return(nil, mockError)

		recorder := httptest.NewRecorder()

		reqBody, _ := json.Marshal(gin.H{
			"refreshToken": validTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())
		mockTokenService.AssertCalled(t, "ValidateRefreshToken", validTokenString)
		mockUserService.AssertCalled(t, "GetUser", getArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", newPairArgs...)
	})

	test.Run("Success", func(t *testing.T) {
		validTokenString := "anothervalid"
		mockTokenID, _ := uuid.NewRandom()
		mockUserID, _ := uuid.NewRandom()

		mockRefreshTokenResp := &models.RefreshToken{
			Token:  validTokenString,
			ID:     mockTokenID,
			UserID: mockUserID,
		}

		mockTokenService.
			On("ValidateRefreshToken", validTokenString).
			Return(mockRefreshTokenResp, nil)

		mockUserResp := &models.User{
			UID: mockUserID,
		}
		getArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockRefreshTokenResp.UserID,
		}

		mockUserService.
			On("GetUser", getArgs...).
			Return(mockUserResp, nil)

		mockNewTokenID, _ := uuid.NewRandom()
		mockNewUserID, _ := uuid.NewRandom()
		mockTokenPairResp := &models.TokenPair{
			AccessToken: models.AccessToken{Token: "aNewIDToken"},
			RefreshToken: models.RefreshToken{
				Token:  "aNewRefreshToken",
				ID:     mockNewTokenID,
				UserID: mockNewUserID,
			},
		}

		newPairArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUserResp,
			mockRefreshTokenResp.ID.String(),
		}

		mockTokenService.
			On("NewPairFromUser", newPairArgs...).
			Return(mockTokenPairResp, nil)

		recorder := httptest.NewRecorder()

		reqBody, _ := json.Marshal(gin.H{
			"refreshToken": validTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)

		respBody, _ := json.Marshal(gin.H{
			"tokens": mockTokenPairResp,
		})

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, respBody, recorder.Body.Bytes())
		mockTokenService.AssertCalled(t, "ValidateRefreshToken", validTokenString)
		mockUserService.AssertCalled(t, "GetUser", getArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", newPairArgs...)
	})

}
