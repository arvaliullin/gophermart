package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderHandler_Submit(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		body           string
		setup          func(*mocks.MockOrderService)
		wantStatusCode int
	}{
		{
			name:   "success - new order",
			userID: 1,
			body:   "12345678903",
			setup: func(orderService *mocks.MockOrderService) {
				orderService.EXPECT().
					SubmitOrder(gomock.Any(), int64(1), "12345678903").
					Return(false, nil)
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name:   "order already exists for same user",
			userID: 1,
			body:   "12345678903",
			setup: func(orderService *mocks.MockOrderService) {
				orderService.EXPECT().
					SubmitOrder(gomock.Any(), int64(1), "12345678903").
					Return(true, nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:   "order belongs to another user",
			userID: 1,
			body:   "12345678903",
			setup: func(orderService *mocks.MockOrderService) {
				orderService.EXPECT().
					SubmitOrder(gomock.Any(), int64(1), "12345678903").
					Return(false, domain.ErrOrderBelongsToOther)
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name:   "invalid luhn number",
			userID: 1,
			body:   "12345678901",
			setup: func(orderService *mocks.MockOrderService) {
				orderService.EXPECT().
					SubmitOrder(gomock.Any(), int64(1), "12345678901").
					Return(false, domain.ErrInvalidOrderNumber)
			},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty body",
			userID:         1,
			body:           "",
			setup:          func(orderService *mocks.MockOrderService) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "unauthorized",
			userID:         0,
			body:           "12345678903",
			setup:          func(orderService *mocks.MockOrderService) {},
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderService := mocks.NewMockOrderService(ctrl)
			tt.setup(orderService)

			handler := handlers.NewOrderHandler(orderService)

			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain")

			if tt.userID > 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()

			handler.Submit(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
		})
	}
}

func TestOrderHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		setup          func(*mocks.MockOrderService)
		wantStatusCode int
		wantBody       bool
	}{
		{
			name:   "success with orders",
			userID: 1,
			setup: func(orderService *mocks.MockOrderService) {
				accrual := decimal.NewFromFloat(500.0)
				orderService.EXPECT().
					GetUserOrders(gomock.Any(), int64(1)).
					Return([]*domain.Order{
						{
							ID:         1,
							UserID:     1,
							Number:     "12345678903",
							Status:     domain.OrderStatusProcessed,
							Accrual:    &accrual,
							UploadedAt: time.Now(),
						},
					}, nil)
			},
			wantStatusCode: http.StatusOK,
			wantBody:       true,
		},
		{
			name:   "no orders",
			userID: 1,
			setup: func(orderService *mocks.MockOrderService) {
				orderService.EXPECT().
					GetUserOrders(gomock.Any(), int64(1)).
					Return([]*domain.Order{}, nil)
			},
			wantStatusCode: http.StatusNoContent,
			wantBody:       false,
		},
		{
			name:           "unauthorized",
			userID:         0,
			setup:          func(orderService *mocks.MockOrderService) {},
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderService := mocks.NewMockOrderService(ctrl)
			tt.setup(orderService)

			handler := handlers.NewOrderHandler(orderService)

			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)

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
