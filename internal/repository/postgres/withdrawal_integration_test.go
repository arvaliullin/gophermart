//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalRepository_GetByUserID(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	balanceRepo := postgres.NewBalanceRepository(testPool)
	withdrawalRepo := postgres.NewWithdrawalRepository(testPool)

	user, err := userRepo.Create(ctx, "withdrawaluser", "password")
	require.NoError(t, err)

	t.Run("пустой список для нового пользователя", func(t *testing.T) {
		withdrawals, err := withdrawalRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Empty(t, withdrawals)
	})

	t.Run("список списаний после операций", func(t *testing.T) {
		err := balanceRepo.AddAccrual(ctx, user.ID, decimal.NewFromFloat(500.0))
		require.NoError(t, err)

		err = balanceRepo.Withdraw(ctx, user.ID, "1111111111", decimal.NewFromFloat(100.0))
		require.NoError(t, err)
		err = balanceRepo.Withdraw(ctx, user.ID, "2222222222", decimal.NewFromFloat(50.0))
		require.NoError(t, err)

		withdrawals, err := withdrawalRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, withdrawals, 2)

		assert.Equal(t, "2222222222", withdrawals[0].OrderNumber)
		assert.True(t, decimal.NewFromFloat(50.0).Equal(withdrawals[0].Sum))
		assert.Equal(t, "1111111111", withdrawals[1].OrderNumber)
		assert.True(t, decimal.NewFromFloat(100.0).Equal(withdrawals[1].Sum))
	})

	t.Run("изоляция списаний между пользователями", func(t *testing.T) {
		user2, err := userRepo.Create(ctx, "anotheruser", "password")
		require.NoError(t, err)

		err = balanceRepo.AddAccrual(ctx, user2.ID, decimal.NewFromFloat(100.0))
		require.NoError(t, err)
		err = balanceRepo.Withdraw(ctx, user2.ID, "3333333333", decimal.NewFromFloat(25.0))
		require.NoError(t, err)

		user1Withdrawals, err := withdrawalRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, user1Withdrawals, 2)

		user2Withdrawals, err := withdrawalRepo.GetByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Len(t, user2Withdrawals, 1)
		assert.Equal(t, "3333333333", user2Withdrawals[0].OrderNumber)
	})
}
