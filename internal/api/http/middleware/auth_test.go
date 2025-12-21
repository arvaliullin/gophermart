package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth_Success_BearerToken(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret")
	token, err := jwtManager.GenerateToken(123)
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())
		assert.Equal(t, int64(123), userID)
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := middleware.Auth(jwtManager)
	protectedHandler := authMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuth_Success_Cookie(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret")
	token, err := jwtManager.GenerateToken(456)
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())
		assert.Equal(t, int64(456), userID)
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := middleware.Auth(jwtManager)
	protectedHandler := authMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuth_NoToken(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := middleware.Auth(jwtManager)
	protectedHandler := authMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuth_InvalidToken(t *testing.T) {
	jwtManager := jwt.NewManager("test-secret")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := middleware.Auth(jwtManager)
	protectedHandler := authMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetUserID_NoUserID(t *testing.T) {
	ctx := context.Background()
	userID := middleware.GetUserID(ctx)
	assert.Equal(t, int64(0), userID)
}

func TestGetUserID_WithUserID(t *testing.T) {
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, int64(789))
	userID := middleware.GetUserID(ctx)
	assert.Equal(t, int64(789), userID)
}
