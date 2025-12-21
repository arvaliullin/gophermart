package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/arvaliullin/gophermart/internal/api/http"
	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewRouter_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockAuthService(ctrl)
	orderService := mocks.NewMockOrderService(ctrl)
	balanceService := mocks.NewMockBalanceService(ctrl)

	jwtManager := jwt.NewManager("test-secret")
	logger := zerolog.Nop()

	router := httpapi.NewRouter(&httpapi.RouterConfig{
		AuthHandler:       handlers.NewAuthHandler(authService),
		OrderHandler:      handlers.NewOrderHandler(orderService),
		BalanceHandler:    handlers.NewBalanceHandler(balanceService),
		WithdrawalHandler: handlers.NewWithdrawalHandler(balanceService),
		JWTManager:        jwtManager,
		Logger:            logger,
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNewRouter_ProtectedRoutes_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockAuthService(ctrl)
	orderService := mocks.NewMockOrderService(ctrl)
	balanceService := mocks.NewMockBalanceService(ctrl)

	jwtManager := jwt.NewManager("test-secret")
	logger := zerolog.Nop()

	router := httpapi.NewRouter(&httpapi.RouterConfig{
		AuthHandler:       handlers.NewAuthHandler(authService),
		OrderHandler:      handlers.NewOrderHandler(orderService),
		BalanceHandler:    handlers.NewBalanceHandler(balanceService),
		WithdrawalHandler: handlers.NewWithdrawalHandler(balanceService),
		JWTManager:        jwtManager,
		Logger:            logger,
	})

	protectedRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/user/orders"},
		{http.MethodGet, "/api/user/orders"},
		{http.MethodGet, "/api/user/balance"},
		{http.MethodPost, "/api/user/balance/withdraw"},
		{http.MethodGet, "/api/user/withdrawals"},
	}

	for _, route := range protectedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
	}
}
