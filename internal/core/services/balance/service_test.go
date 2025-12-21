package balance_test

import (
	"context"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/arvaliullin/gophermart/internal/core/services/balance"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_GetBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	withdrawalRepo := mocks.NewMockWithdrawalRepository(ctrl)
	service := balance.NewService(balanceRepo, withdrawalRepo)

	expectedBalance := &domain.Balance{
		UserID:    1,
		Current:   decimal.NewFromFloat(500.5),
		Withdrawn: decimal.NewFromFloat(100.0),
	}

	balanceRepo.EXPECT().
		GetByUserID(gomock.Any(), int64(1)).
		Return(expectedBalance, nil)

	result, err := service.GetBalance(context.Background(), 1)

	require.NoError(t, err)
	assert.True(t, expectedBalance.Current.Equal(result.Current))
	assert.True(t, expectedBalance.Withdrawn.Equal(result.Withdrawn))
}

func TestService_Withdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	withdrawalRepo := mocks.NewMockWithdrawalRepository(ctrl)
	service := balance.NewService(balanceRepo, withdrawalRepo)

	amount := decimal.NewFromInt(100)
	balanceRepo.EXPECT().
		Withdraw(gomock.Any(), int64(1), "79927398713", amount).
		Return(nil)

	err := service.Withdraw(context.Background(), 1, "79927398713", amount)

	require.NoError(t, err)
}

func TestService_Withdraw_InvalidOrderNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	withdrawalRepo := mocks.NewMockWithdrawalRepository(ctrl)
	service := balance.NewService(balanceRepo, withdrawalRepo)

	err := service.Withdraw(context.Background(), 1, "invalid-number", decimal.NewFromInt(100))

	assert.ErrorIs(t, err, domain.ErrInvalidOrderNumber)
}

func TestService_Withdraw_InsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	withdrawalRepo := mocks.NewMockWithdrawalRepository(ctrl)
	service := balance.NewService(balanceRepo, withdrawalRepo)

	amount := decimal.NewFromInt(1000)
	balanceRepo.EXPECT().
		Withdraw(gomock.Any(), int64(1), "79927398713", amount).
		Return(domain.ErrInsufficientBalance)

	err := service.Withdraw(context.Background(), 1, "79927398713", amount)

	assert.ErrorIs(t, err, domain.ErrInsufficientBalance)
}

func TestService_GetWithdrawals_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	withdrawalRepo := mocks.NewMockWithdrawalRepository(ctrl)
	service := balance.NewService(balanceRepo, withdrawalRepo)

	withdrawals := []*domain.Withdrawal{
		{
			ID:          1,
			UserID:      1,
			OrderNumber: "12345678903",
			Sum:         decimal.NewFromInt(100),
			ProcessedAt: time.Now(),
		},
		{
			ID:          2,
			UserID:      1,
			OrderNumber: "79927398713",
			Sum:         decimal.NewFromInt(50),
			ProcessedAt: time.Now(),
		},
	}

	withdrawalRepo.EXPECT().
		GetByUserID(gomock.Any(), int64(1)).
		Return(withdrawals, nil)

	result, err := service.GetWithdrawals(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestService_GetWithdrawals_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	withdrawalRepo := mocks.NewMockWithdrawalRepository(ctrl)
	service := balance.NewService(balanceRepo, withdrawalRepo)

	withdrawalRepo.EXPECT().
		GetByUserID(gomock.Any(), int64(1)).
		Return([]*domain.Withdrawal{}, nil)

	result, err := service.GetWithdrawals(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, result)
}
