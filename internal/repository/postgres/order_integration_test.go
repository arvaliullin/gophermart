//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderRepository_Create(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	orderRepo := postgres.NewOrderRepository(testPool)

	user, err := userRepo.Create(ctx, "orderuser", "password")
	require.NoError(t, err)

	t.Run("успешное создание заказа", func(t *testing.T) {
		order, err := orderRepo.Create(ctx, user.ID, "12345678903")
		require.NoError(t, err)
		assert.NotZero(t, order.ID)
		assert.Equal(t, user.ID, order.UserID)
		assert.Equal(t, "12345678903", order.Number)
		assert.Equal(t, domain.OrderStatusNew, order.Status)
		assert.Nil(t, order.Accrual)
		assert.NotZero(t, order.UploadedAt)
	})

	t.Run("ошибка при дублировании номера заказа", func(t *testing.T) {
		_, err := orderRepo.Create(ctx, user.ID, "9876543210")
		require.NoError(t, err)

		_, err = orderRepo.Create(ctx, user.ID, "9876543210")
		assert.ErrorIs(t, err, domain.ErrOrderAlreadyExists)
	})
}

func TestOrderRepository_GetByNumber(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	orderRepo := postgres.NewOrderRepository(testPool)

	user, err := userRepo.Create(ctx, "getbynum", "password")
	require.NoError(t, err)

	t.Run("успешное получение заказа по номеру", func(t *testing.T) {
		created, err := orderRepo.Create(ctx, user.ID, "1111111111")
		require.NoError(t, err)

		found, err := orderRepo.GetByNumber(ctx, "1111111111")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, created.Number, found.Number)
	})

	t.Run("ошибка при несуществующем номере", func(t *testing.T) {
		_, err := orderRepo.GetByNumber(ctx, "0000000000")
		assert.ErrorIs(t, err, domain.ErrOrderNotFound)
	})
}

func TestOrderRepository_GetByUserID(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	orderRepo := postgres.NewOrderRepository(testPool)

	user, err := userRepo.Create(ctx, "getbyuser", "password")
	require.NoError(t, err)

	t.Run("успешное получение заказов пользователя", func(t *testing.T) {
		_, err := orderRepo.Create(ctx, user.ID, "2222222222")
		require.NoError(t, err)
		_, err = orderRepo.Create(ctx, user.ID, "3333333333")
		require.NoError(t, err)

		orders, err := orderRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, orders, 2)
		assert.Equal(t, "3333333333", orders[0].Number)
		assert.Equal(t, "2222222222", orders[1].Number)
	})

	t.Run("пустой список для нового пользователя", func(t *testing.T) {
		user2, err := userRepo.Create(ctx, "noorders", "password")
		require.NoError(t, err)

		orders, err := orderRepo.GetByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Empty(t, orders)
	})
}

func TestOrderRepository_GetPendingOrders(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	orderRepo := postgres.NewOrderRepository(testPool)

	user, err := userRepo.Create(ctx, "pending", "password")
	require.NoError(t, err)

	t.Run("получение необработанных заказов", func(t *testing.T) {
		_, err := orderRepo.Create(ctx, user.ID, "4444444444")
		require.NoError(t, err)
		_, err = orderRepo.Create(ctx, user.ID, "5555555555")
		require.NoError(t, err)

		accrual := decimal.NewFromFloat(100.0)
		err = orderRepo.UpdateStatus(ctx, "5555555555", domain.OrderStatusProcessed, &accrual)
		require.NoError(t, err)

		pending, err := orderRepo.GetPendingOrders(ctx)
		require.NoError(t, err)
		assert.Len(t, pending, 1)
		assert.Equal(t, "4444444444", pending[0].Number)
	})
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	userRepo := postgres.NewUserRepository(testPool)
	orderRepo := postgres.NewOrderRepository(testPool)

	user, err := userRepo.Create(ctx, "updatestat", "password")
	require.NoError(t, err)

	t.Run("успешное обновление статуса", func(t *testing.T) {
		_, err := orderRepo.Create(ctx, user.ID, "6666666666")
		require.NoError(t, err)

		accrual := decimal.NewFromFloat(250.50)
		err = orderRepo.UpdateStatus(ctx, "6666666666", domain.OrderStatusProcessed, &accrual)
		require.NoError(t, err)

		updated, err := orderRepo.GetByNumber(ctx, "6666666666")
		require.NoError(t, err)
		assert.Equal(t, domain.OrderStatusProcessed, updated.Status)
		require.NotNil(t, updated.Accrual)
		assert.True(t, decimal.NewFromFloat(250.50).Equal(*updated.Accrual))
	})

	t.Run("обновление статуса на INVALID без начисления", func(t *testing.T) {
		_, err := orderRepo.Create(ctx, user.ID, "7777777777")
		require.NoError(t, err)

		err = orderRepo.UpdateStatus(ctx, "7777777777", domain.OrderStatusInvalid, nil)
		require.NoError(t, err)

		updated, err := orderRepo.GetByNumber(ctx, "7777777777")
		require.NoError(t, err)
		assert.Equal(t, domain.OrderStatusInvalid, updated.Status)
		assert.Nil(t, updated.Accrual)
	})

	t.Run("ошибка при обновлении несуществующего заказа", func(t *testing.T) {
		err := orderRepo.UpdateStatus(ctx, "0000000000", domain.OrderStatusProcessed, nil)
		assert.ErrorIs(t, err, domain.ErrOrderNotFound)
	})
}
