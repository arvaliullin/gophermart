package accrual_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorker_ProcessOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	accrualClient := mocks.NewMockAccrualClient(ctrl)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	worker := accrual.NewWorker(orderRepo, balanceRepo, accrualClient, logger)

	pendingOrders := []*domain.Order{
		{
			ID:     1,
			UserID: 1,
			Number: "12345678903",
			Status: domain.OrderStatusNew,
		},
	}

	orderRepo.EXPECT().
		GetPendingOrders(gomock.Any()).
		Return(pendingOrders, nil).
		AnyTimes()

	accrualClient.EXPECT().
		GetOrderAccrual(gomock.Any(), "12345678903").
		Return(&ports.AccrualResponse{
			Order:   "12345678903",
			Status:  domain.OrderStatusProcessed,
			Accrual: 500.0,
		}, nil).
		AnyTimes()

	orderRepo.EXPECT().
		UpdateStatus(gomock.Any(), "12345678903", domain.OrderStatusProcessed, gomock.Any()).
		Return(nil).
		AnyTimes()

	balanceRepo.EXPECT().
		AddAccrual(gomock.Any(), int64(1), 500.0).
		Return(nil).
		AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	go worker.Run(ctx)
	<-ctx.Done()
}

func TestWorker_ProcessOrder_NoContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	accrualClient := mocks.NewMockAccrualClient(ctrl)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	worker := accrual.NewWorker(orderRepo, balanceRepo, accrualClient, logger)

	pendingOrders := []*domain.Order{
		{
			ID:     1,
			UserID: 1,
			Number: "12345678903",
			Status: domain.OrderStatusNew,
		},
	}

	orderRepo.EXPECT().
		GetPendingOrders(gomock.Any()).
		Return(pendingOrders, nil).
		AnyTimes()

	accrualClient.EXPECT().
		GetOrderAccrual(gomock.Any(), "12345678903").
		Return(nil, nil).
		AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	go worker.Run(ctx)
	<-ctx.Done()
}

func TestWorker_ProcessOrder_Processing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	accrualClient := mocks.NewMockAccrualClient(ctrl)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	worker := accrual.NewWorker(orderRepo, balanceRepo, accrualClient, logger)

	pendingOrders := []*domain.Order{
		{
			ID:     1,
			UserID: 1,
			Number: "12345678903",
			Status: domain.OrderStatusNew,
		},
	}

	orderRepo.EXPECT().
		GetPendingOrders(gomock.Any()).
		Return(pendingOrders, nil).
		AnyTimes()

	accrualClient.EXPECT().
		GetOrderAccrual(gomock.Any(), "12345678903").
		Return(&ports.AccrualResponse{
			Order:  "12345678903",
			Status: domain.OrderStatusProcessing,
		}, nil).
		AnyTimes()

	orderRepo.EXPECT().
		UpdateStatus(gomock.Any(), "12345678903", domain.OrderStatusProcessing, gomock.Any()).
		Return(nil).
		AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	go worker.Run(ctx)
	<-ctx.Done()
}

func TestWorker_ProcessOrder_RateLimited(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	accrualClient := mocks.NewMockAccrualClient(ctrl)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	worker := accrual.NewWorker(orderRepo, balanceRepo, accrualClient, logger)

	pendingOrders := []*domain.Order{
		{
			ID:     1,
			UserID: 1,
			Number: "12345678903",
			Status: domain.OrderStatusNew,
		},
	}

	orderRepo.EXPECT().
		GetPendingOrders(gomock.Any()).
		Return(pendingOrders, nil).
		AnyTimes()

	accrualClient.EXPECT().
		GetOrderAccrual(gomock.Any(), "12345678903").
		Return(nil, &accrual.RetryAfterError{Duration: 60 * time.Second}).
		AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	go worker.Run(ctx)
	<-ctx.Done()
}

func TestWorker_NoPendingOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	accrualClient := mocks.NewMockAccrualClient(ctrl)
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	worker := accrual.NewWorker(orderRepo, balanceRepo, accrualClient, logger)

	orderRepo.EXPECT().
		GetPendingOrders(gomock.Any()).
		Return([]*domain.Order{}, nil).
		AnyTimes()

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	go worker.Run(ctx)
	<-ctx.Done()
}

func TestRetryAfterError_Error(t *testing.T) {
	err := &accrual.RetryAfterError{Duration: 60 * time.Second}
	assert.Equal(t, "превышен лимит запросов, повторить через 1m0s", err.Error())
}

