package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setup          func(*mocks.MockAuthService)
		wantStatusCode int
		wantCookie     bool
	}{
		{
			name: "success",
			body: map[string]string{"login": "testuser", "password": "password123"},
			setup: func(authService *mocks.MockAuthService) {
				authService.EXPECT().
					Register(gomock.Any(), "testuser", "password123").
					Return("valid-jwt-token", nil)
			},
			wantStatusCode: http.StatusOK,
			wantCookie:     true,
		},
		{
			name: "user already exists",
			body: map[string]string{"login": "existinguser", "password": "password123"},
			setup: func(authService *mocks.MockAuthService) {
				authService.EXPECT().
					Register(gomock.Any(), "existinguser", "password123").
					Return("", domain.ErrUserAlreadyExists)
			},
			wantStatusCode: http.StatusConflict,
			wantCookie:     false,
		},
		{
			name:           "invalid json",
			body:           "invalid json",
			setup:          func(authService *mocks.MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantCookie:     false,
		},
		{
			name:           "empty login",
			body:           map[string]string{"login": "", "password": "password123"},
			setup:          func(authService *mocks.MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantCookie:     false,
		},
		{
			name:           "empty password",
			body:           map[string]string{"login": "testuser", "password": ""},
			setup:          func(authService *mocks.MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantCookie:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := mocks.NewMockAuthService(ctrl)
			tt.setup(authService)

			handler := handlers.NewAuthHandler(authService)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Register(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			if tt.wantCookie {
				cookies := rr.Result().Cookies()
				found := false
				for _, c := range cookies {
					if c.Name == "auth_token" {
						found = true
						assert.NotEmpty(t, c.Value)
					}
				}
				assert.True(t, found, "expected auth_token cookie")
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setup          func(*mocks.MockAuthService)
		wantStatusCode int
		wantCookie     bool
	}{
		{
			name: "success",
			body: map[string]string{"login": "testuser", "password": "password123"},
			setup: func(authService *mocks.MockAuthService) {
				authService.EXPECT().
					Login(gomock.Any(), "testuser", "password123").
					Return("valid-jwt-token", nil)
			},
			wantStatusCode: http.StatusOK,
			wantCookie:     true,
		},
		{
			name: "invalid credentials",
			body: map[string]string{"login": "unknown", "password": "password123"},
			setup: func(authService *mocks.MockAuthService) {
				authService.EXPECT().
					Login(gomock.Any(), "unknown", "password123").
					Return("", domain.ErrInvalidCredentials)
			},
			wantStatusCode: http.StatusUnauthorized,
			wantCookie:     false,
		},
		{
			name:           "invalid json",
			body:           "invalid",
			setup:          func(authService *mocks.MockAuthService) {},
			wantStatusCode: http.StatusBadRequest,
			wantCookie:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := mocks.NewMockAuthService(ctrl)
			tt.setup(authService)

			handler := handlers.NewAuthHandler(authService)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Login(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
		})
	}
}
