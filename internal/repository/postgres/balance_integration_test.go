//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceRepository_GetByUserID(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	balanceRepo := postgres.NewBalanceRepository(testPool)

	user, err := userRepo.Create(ctx, "balanceuser", "password")
	require.NoError(t, err)

	t.Run("возвращает нулевой баланс для нового пользователя", func(t *testing.T) {
		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, balance.UserID)
		assert.Equal(t, 0.0, balance.Current)
		assert.Equal(t, 0.0, balance.Withdrawn)
	})

	t.Run("возвращает корректный баланс после создания", func(t *testing.T) {
		err := balanceRepo.CreateForUser(ctx, user.ID)
		require.NoError(t, err)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 0.0, balance.Current)
		assert.Equal(t, 0.0, balance.Withdrawn)
	})
}

func TestBalanceRepository_CreateForUser(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	balanceRepo := postgres.NewBalanceRepository(testPool)

	user, err := userRepo.Create(ctx, "createbal", "password")
	require.NoError(t, err)

	t.Run("успешное создание баланса", func(t *testing.T) {
		err := balanceRepo.CreateForUser(ctx, user.ID)
		require.NoError(t, err)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, balance.UserID)
	})

	t.Run("повторное создание не вызывает ошибку", func(t *testing.T) {
		err := balanceRepo.CreateForUser(ctx, user.ID)
		require.NoError(t, err)
	})
}

func TestBalanceRepository_AddAccrual(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	balanceRepo := postgres.NewBalanceRepository(testPool)

	user, err := userRepo.Create(ctx, "accrualuser", "password")
	require.NoError(t, err)

	t.Run("успешное добавление начисления", func(t *testing.T) {
		err := balanceRepo.AddAccrual(ctx, user.ID, 100.50)
		require.NoError(t, err)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 100.50, balance.Current)
	})

	t.Run("накопление начислений", func(t *testing.T) {
		err := balanceRepo.AddAccrual(ctx, user.ID, 50.25)
		require.NoError(t, err)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 150.75, balance.Current)
	})
}

func TestBalanceRepository_Withdraw(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	balanceRepo := postgres.NewBalanceRepository(testPool)

	user, err := userRepo.Create(ctx, "withdrawuser", "password")
	require.NoError(t, err)

	err = balanceRepo.AddAccrual(ctx, user.ID, 500.0)
	require.NoError(t, err)

	t.Run("успешное списание средств", func(t *testing.T) {
		err := balanceRepo.Withdraw(ctx, user.ID, "2377225624", 100.0)
		require.NoError(t, err)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 400.0, balance.Current)
		assert.Equal(t, 100.0, balance.Withdrawn)
	})

	t.Run("повторное списание", func(t *testing.T) {
		err := balanceRepo.Withdraw(ctx, user.ID, "1234567890", 150.0)
		require.NoError(t, err)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 250.0, balance.Current)
		assert.Equal(t, 250.0, balance.Withdrawn)
	})

	t.Run("ошибка при недостаточном балансе", func(t *testing.T) {
		err := balanceRepo.Withdraw(ctx, user.ID, "9999999999", 1000.0)
		assert.ErrorIs(t, err, domain.ErrInsufficientBalance)

		balance, err := balanceRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, 250.0, balance.Current)
	})

	t.Run("ошибка при списании без баланса", func(t *testing.T) {
		user2, err := userRepo.Create(ctx, "nobalance", "password")
		require.NoError(t, err)

		err = balanceRepo.Withdraw(ctx, user2.ID, "0000000000", 10.0)
		assert.ErrorIs(t, err, domain.ErrInsufficientBalance)
	})
}
