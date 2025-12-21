package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/arvaliullin/gophermart/internal/core/services/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_SubmitOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	service := order.NewService(orderRepo)

	newOrder := &domain.Order{
		ID:         1,
		UserID:     1,
		Number:     "12345678903",
		Status:     domain.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	orderRepo.EXPECT().
		Create(gomock.Any(), int64(1), "12345678903").
		Return(newOrder, nil)

	alreadyExists, err := service.SubmitOrder(context.Background(), 1, "12345678903")

	require.NoError(t, err)
	assert.False(t, alreadyExists)
}

func TestService_SubmitOrder_InvalidLuhn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	service := order.NewService(orderRepo)

	_, err := service.SubmitOrder(context.Background(), 1, "12345678901")

	assert.ErrorIs(t, err, domain.ErrInvalidOrderNumber)
}

func TestService_SubmitOrder_AlreadyExistsSameUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	service := order.NewService(orderRepo)

	existingOrder := &domain.Order{
		ID:         1,
		UserID:     1,
		Number:     "12345678903",
		Status:     domain.OrderStatusProcessed,
		UploadedAt: time.Now(),
	}

	orderRepo.EXPECT().
		Create(gomock.Any(), int64(1), "12345678903").
		Return(nil, domain.ErrOrderAlreadyExists)

	orderRepo.EXPECT().
		GetByNumber(gomock.Any(), "12345678903").
		Return(existingOrder, nil)

	alreadyExists, err := service.SubmitOrder(context.Background(), 1, "12345678903")

	require.NoError(t, err)
	assert.True(t, alreadyExists)
}

func TestService_SubmitOrder_BelongsToOther(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	service := order.NewService(orderRepo)

	existingOrder := &domain.Order{
		ID:         1,
		UserID:     2,
		Number:     "12345678903",
		Status:     domain.OrderStatusProcessed,
		UploadedAt: time.Now(),
	}

	orderRepo.EXPECT().
		Create(gomock.Any(), int64(1), "12345678903").
		Return(nil, domain.ErrOrderAlreadyExists)

	orderRepo.EXPECT().
		GetByNumber(gomock.Any(), "12345678903").
		Return(existingOrder, nil)

	_, err := service.SubmitOrder(context.Background(), 1, "12345678903")

	assert.ErrorIs(t, err, domain.ErrOrderBelongsToOther)
}

func TestService_GetUserOrders_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	service := order.NewService(orderRepo)

	accrual := 500.0
	orders := []*domain.Order{
		{
			ID:         1,
			UserID:     1,
			Number:     "12345678903",
			Status:     domain.OrderStatusProcessed,
			Accrual:    &accrual,
			UploadedAt: time.Now(),
		},
		{
			ID:         2,
			UserID:     1,
			Number:     "79927398713",
			Status:     domain.OrderStatusNew,
			UploadedAt: time.Now(),
		},
	}

	orderRepo.EXPECT().
		GetByUserID(gomock.Any(), int64(1)).
		Return(orders, nil)

	result, err := service.GetUserOrders(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestService_GetUserOrders_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	service := order.NewService(orderRepo)

	orderRepo.EXPECT().
		GetByUserID(gomock.Any(), int64(1)).
		Return([]*domain.Order{}, nil)

	result, err := service.GetUserOrders(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, result)
}

