package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func testStrategy() *retry.Strategy {
	return retry.NewStrategy(
		[]time.Duration{10 * time.Millisecond},
		func(err error) bool { return false },
	)
}

func TestNewUserRepositoryAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("успешное создание", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(ctrl)
		adapter, err := NewUserRepositoryAdapter(repo, testStrategy())
		require.NoError(t, err)
		assert.NotNil(t, adapter)
	})

	t.Run("ошибка при nil репозитории", func(t *testing.T) {
		adapter, err := NewUserRepositoryAdapter(nil, testStrategy())
		assert.ErrorIs(t, err, ErrUserRepoNil)
		assert.Nil(t, adapter)
	})
}

func TestUserRepositoryAdapter_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockUserRepository(ctrl)
	adapter, _ := NewUserRepositoryAdapter(repo, testStrategy())

	expectedUser := &domain.User{ID: 1, Login: "test"}
	repo.EXPECT().Create(ctx, "test", "hash").Return(expectedUser, nil)

	user, err := adapter.Create(ctx, "test", "hash")
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserRepositoryAdapter_GetByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockUserRepository(ctrl)
	adapter, _ := NewUserRepositoryAdapter(repo, testStrategy())

	expectedUser := &domain.User{ID: 1, Login: "test"}
	repo.EXPECT().GetByLogin(ctx, "test").Return(expectedUser, nil)

	user, err := adapter.GetByLogin(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserRepositoryAdapter_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockUserRepository(ctrl)
	adapter, _ := NewUserRepositoryAdapter(repo, testStrategy())

	expectedUser := &domain.User{ID: 1, Login: "test"}
	repo.EXPECT().GetByID(ctx, int64(1)).Return(expectedUser, nil)

	user, err := adapter.GetByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestNewOrderRepositoryAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("успешное создание", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(ctrl)
		adapter, err := NewOrderRepositoryAdapter(repo, testStrategy())
		require.NoError(t, err)
		assert.NotNil(t, adapter)
	})

	t.Run("ошибка при nil репозитории", func(t *testing.T) {
		adapter, err := NewOrderRepositoryAdapter(nil, testStrategy())
		assert.ErrorIs(t, err, ErrOrderRepoNil)
		assert.Nil(t, adapter)
	})
}

func TestOrderRepositoryAdapter_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockOrderRepository(ctrl)
	adapter, _ := NewOrderRepositoryAdapter(repo, testStrategy())

	expectedOrder := &domain.Order{ID: 1, Number: "123"}
	repo.EXPECT().Create(ctx, int64(1), "123").Return(expectedOrder, nil)

	order, err := adapter.Create(ctx, 1, "123")
	require.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestOrderRepositoryAdapter_GetByNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockOrderRepository(ctrl)
	adapter, _ := NewOrderRepositoryAdapter(repo, testStrategy())

	expectedOrder := &domain.Order{ID: 1, Number: "123"}
	repo.EXPECT().GetByNumber(ctx, "123").Return(expectedOrder, nil)

	order, err := adapter.GetByNumber(ctx, "123")
	require.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestOrderRepositoryAdapter_GetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockOrderRepository(ctrl)
	adapter, _ := NewOrderRepositoryAdapter(repo, testStrategy())

	expectedOrders := []*domain.Order{{ID: 1}, {ID: 2}}
	repo.EXPECT().GetByUserID(ctx, int64(1)).Return(expectedOrders, nil)

	orders, err := adapter.GetByUserID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
}

func TestOrderRepositoryAdapter_GetPendingOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockOrderRepository(ctrl)
	adapter, _ := NewOrderRepositoryAdapter(repo, testStrategy())

	expectedOrders := []*domain.Order{{ID: 1, Status: domain.OrderStatusNew}}
	repo.EXPECT().GetPendingOrders(ctx).Return(expectedOrders, nil)

	orders, err := adapter.GetPendingOrders(ctx)
	require.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
}

func TestOrderRepositoryAdapter_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockOrderRepository(ctrl)
	adapter, _ := NewOrderRepositoryAdapter(repo, testStrategy())

	accrual := decimal.NewFromFloat(100.0)
	repo.EXPECT().UpdateStatus(ctx, "123", domain.OrderStatusProcessed, &accrual).Return(nil)

	err := adapter.UpdateStatus(ctx, "123", domain.OrderStatusProcessed, &accrual)
	require.NoError(t, err)
}

func TestNewBalanceRepositoryAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("успешное создание", func(t *testing.T) {
		repo := mocks.NewMockBalanceRepository(ctrl)
		adapter, err := NewBalanceRepositoryAdapter(repo, testStrategy())
		require.NoError(t, err)
		assert.NotNil(t, adapter)
	})

	t.Run("ошибка при nil репозитории", func(t *testing.T) {
		adapter, err := NewBalanceRepositoryAdapter(nil, testStrategy())
		assert.ErrorIs(t, err, ErrBalanceRepoNil)
		assert.Nil(t, adapter)
	})
}

func TestBalanceRepositoryAdapter_GetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockBalanceRepository(ctrl)
	adapter, _ := NewBalanceRepositoryAdapter(repo, testStrategy())

	expectedBalance := &domain.Balance{UserID: 1, Current: decimal.NewFromInt(100)}
	repo.EXPECT().GetByUserID(ctx, int64(1)).Return(expectedBalance, nil)

	balance, err := adapter.GetByUserID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

func TestBalanceRepositoryAdapter_CreateForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockBalanceRepository(ctrl)
	adapter, _ := NewBalanceRepositoryAdapter(repo, testStrategy())

	repo.EXPECT().CreateForUser(ctx, int64(1)).Return(nil)

	err := adapter.CreateForUser(ctx, 1)
	require.NoError(t, err)
}

func TestBalanceRepositoryAdapter_AddAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockBalanceRepository(ctrl)
	adapter, _ := NewBalanceRepositoryAdapter(repo, testStrategy())

	amount := decimal.NewFromFloat(50.0)
	repo.EXPECT().AddAccrual(ctx, int64(1), amount).Return(nil)

	err := adapter.AddAccrual(ctx, 1, amount)
	require.NoError(t, err)
}

func TestBalanceRepositoryAdapter_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockBalanceRepository(ctrl)
	adapter, _ := NewBalanceRepositoryAdapter(repo, testStrategy())

	amount := decimal.NewFromFloat(25.0)
	repo.EXPECT().Withdraw(ctx, int64(1), "123", amount).Return(nil)

	err := adapter.Withdraw(ctx, 1, "123", amount)
	require.NoError(t, err)
}

func TestNewWithdrawalRepositoryAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("успешное создание", func(t *testing.T) {
		repo := mocks.NewMockWithdrawalRepository(ctrl)
		adapter, err := NewWithdrawalRepositoryAdapter(repo, testStrategy())
		require.NoError(t, err)
		assert.NotNil(t, adapter)
	})

	t.Run("ошибка при nil репозитории", func(t *testing.T) {
		adapter, err := NewWithdrawalRepositoryAdapter(nil, testStrategy())
		assert.ErrorIs(t, err, ErrWithdrawalRepoNil)
		assert.Nil(t, adapter)
	})
}

func TestWithdrawalRepositoryAdapter_GetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockWithdrawalRepository(ctrl)
	adapter, _ := NewWithdrawalRepositoryAdapter(repo, testStrategy())

	expectedWithdrawals := []*domain.Withdrawal{{ID: 1}, {ID: 2}}
	repo.EXPECT().GetByUserID(ctx, int64(1)).Return(expectedWithdrawals, nil)

	withdrawals, err := adapter.GetByUserID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedWithdrawals, withdrawals)
}

func TestRetryOnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockUserRepository(ctrl)

	retryableErr := errors.New("connection error")
	strategy := retry.NewStrategy(
		[]time.Duration{10 * time.Millisecond, 20 * time.Millisecond},
		func(err error) bool { return err == retryableErr },
	)

	adapter, _ := NewUserRepositoryAdapter(repo, strategy)

	expectedUser := &domain.User{ID: 1, Login: "test"}
	gomock.InOrder(
		repo.EXPECT().GetByLogin(ctx, "test").Return(nil, retryableErr),
		repo.EXPECT().GetByLogin(ctx, "test").Return(expectedUser, nil),
	)

	user, err := adapter.GetByLogin(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}
