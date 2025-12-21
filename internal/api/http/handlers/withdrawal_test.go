package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWithdrawalHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		setup          func(*mocks.MockBalanceService)
		wantStatusCode int
		wantBody       bool
	}{
		{
			name:   "success with withdrawals",
			userID: 1,
			setup: func(balanceService *mocks.MockBalanceService) {
				balanceService.EXPECT().
					GetWithdrawals(gomock.Any(), int64(1)).
					Return([]*domain.Withdrawal{
						{
							ID:          1,
							UserID:      1,
							OrderNumber: "12345678903",
							Sum:         100,
							ProcessedAt: time.Now(),
						},
						{
							ID:          2,
							UserID:      1,
							OrderNumber: "79927398713",
							Sum:         50,
							ProcessedAt: time.Now(),
						},
					}, nil)
			},
			wantStatusCode: http.StatusOK,
			wantBody:       true,
		},
		{
			name:   "no withdrawals",
			userID: 1,
			setup: func(balanceService *mocks.MockBalanceService) {
				balanceService.EXPECT().
					GetWithdrawals(gomock.Any(), int64(1)).
					Return([]*domain.Withdrawal{}, nil)
			},
			wantStatusCode: http.StatusNoContent,
			wantBody:       false,
		},
		{
			name:           "unauthorized",
			userID:         0,
			setup:          func(balanceService *mocks.MockBalanceService) {},
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			balanceService := mocks.NewMockBalanceService(ctrl)
			tt.setup(balanceService)

			handler := handlers.NewWithdrawalHandler(balanceService)

			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)

			if tt.userID > 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()

			handler.List(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			if tt.wantBody {
				assert.NotEmpty(t, rr.Body.String())
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			}
		})
	}
}
