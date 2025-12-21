package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBalanceHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		setup          func(*mocks.MockBalanceService)
		wantStatusCode int
		wantCurrent    string
		wantWithdrawn  string
	}{
		{
			name:   "success",
			userID: 1,
			setup: func(balanceService *mocks.MockBalanceService) {
				balanceService.EXPECT().
					GetBalance(gomock.Any(), int64(1)).
					Return(&domain.Balance{
						UserID:    1,
						Current:   decimal.NewFromFloat(500.5),
						Withdrawn: decimal.NewFromFloat(100.0),
					}, nil)
			},
			wantStatusCode: http.StatusOK,
			wantCurrent:    "500.5",
			wantWithdrawn:  "100",
		},
		{
			name:           "unauthorized",
			userID:         0,
			setup:          func(balanceService *mocks.MockBalanceService) {},
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			balanceService := mocks.NewMockBalanceService(ctrl)
			tt.setup(balanceService)

			handler := handlers.NewBalanceHandler(balanceService)

			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)

			if tt.userID > 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()

			handler.Get(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.wantStatusCode == http.StatusOK {
				var resp map[string]json.Number
				err := json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCurrent, resp["current"].String())
				assert.Equal(t, tt.wantWithdrawn, resp["withdrawn"].String())
			}
		})
	}
}

func TestBalanceHandler_Withdraw(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		body           any
		setup          func(*mocks.MockBalanceService)
		wantStatusCode int
	}{
		{
			name:   "success",
			userID: 1,
			body:   map[string]any{"order": "79927398713", "sum": 100.0},
			setup: func(balanceService *mocks.MockBalanceService) {
				balanceService.EXPECT().
					Withdraw(gomock.Any(), int64(1), "79927398713", gomock.Any()).
					Return(nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:   "insufficient balance",
			userID: 1,
			body:   map[string]any{"order": "79927398713", "sum": 1000.0},
			setup: func(balanceService *mocks.MockBalanceService) {
				balanceService.EXPECT().
					Withdraw(gomock.Any(), int64(1), "79927398713", gomock.Any()).
					Return(domain.ErrInsufficientBalance)
			},
			wantStatusCode: http.StatusPaymentRequired,
		},
		{
			name:   "invalid order number",
			userID: 1,
			body:   map[string]any{"order": "12345", "sum": 100.0},
			setup: func(balanceService *mocks.MockBalanceService) {
				balanceService.EXPECT().
					Withdraw(gomock.Any(), int64(1), "12345", gomock.Any()).
					Return(domain.ErrInvalidOrderNumber)
			},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid json",
			userID:         1,
			body:           "invalid",
			setup:          func(balanceService *mocks.MockBalanceService) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "unauthorized",
			userID:         0,
			body:           map[string]any{"order": "79927398713", "sum": 100.0},
			setup:          func(balanceService *mocks.MockBalanceService) {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "empty order",
			userID:         1,
			body:           map[string]any{"order": "", "sum": 100.0},
			setup:          func(balanceService *mocks.MockBalanceService) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "zero sum",
			userID:         1,
			body:           map[string]any{"order": "79927398713", "sum": 0},
			setup:          func(balanceService *mocks.MockBalanceService) {},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			balanceService := mocks.NewMockBalanceService(ctrl)
			tt.setup(balanceService)

			handler := handlers.NewBalanceHandler(balanceService)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.userID > 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()

			handler.Withdraw(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
		})
	}
}
