package retry

import (
	"context"
	"fmt"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
)

// ErrOrderRepoNil возвращается при попытке создать адаптер с nil репозиторием.
var ErrOrderRepoNil = fmt.Errorf("репозиторий заказов не задан")

// OrderRepositoryAdapter добавляет стратегию повторов для репозитория заказов.
type OrderRepositoryAdapter struct {
	repo     ports.OrderRepository
	strategy *retry.Strategy
}

// NewOrderRepositoryAdapter создаёт адаптер репозитория заказов с поддержкой retry.
func NewOrderRepositoryAdapter(repo ports.OrderRepository, strategy *retry.Strategy) (*OrderRepositoryAdapter, error) {
	if repo == nil {
		return nil, ErrOrderRepoNil
	}

	return &OrderRepositoryAdapter{
		repo:     repo,
		strategy: strategy,
	}, nil
}

// Create создаёт новый заказ.
func (a *OrderRepositoryAdapter) Create(ctx context.Context, userID int64, number string) (*domain.Order, error) {
	var order *domain.Order
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		order, err = a.repo.Create(ctx, userID, number)
		return err
	})
	return order, err
}

// GetByNumber возвращает заказ по номеру.
func (a *OrderRepositoryAdapter) GetByNumber(ctx context.Context, number string) (*domain.Order, error) {
	var order *domain.Order
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		order, err = a.repo.GetByNumber(ctx, number)
		return err
	})
	return order, err
}

// GetByUserID возвращает все заказы пользователя.
func (a *OrderRepositoryAdapter) GetByUserID(ctx context.Context, userID int64) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		orders, err = a.repo.GetByUserID(ctx, userID)
		return err
	})
	return orders, err
}

// GetPendingOrders возвращает заказы со статусами NEW или PROCESSING.
func (a *OrderRepositoryAdapter) GetPendingOrders(ctx context.Context) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		orders, err = a.repo.GetPendingOrders(ctx)
		return err
	})
	return orders, err
}

// UpdateStatus обновляет статус и начисление заказа.
func (a *OrderRepositoryAdapter) UpdateStatus(ctx context.Context, number string, status domain.OrderStatus, accrual *float64) error {
	return a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		return a.repo.UpdateStatus(ctx, number, status, accrual)
	})
}
