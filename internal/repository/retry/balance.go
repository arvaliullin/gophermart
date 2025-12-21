package retry

import (
	"context"
	"fmt"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
	"github.com/shopspring/decimal"
)

// ErrBalanceRepoNil возвращается при попытке создать адаптер с nil репозиторием.
var ErrBalanceRepoNil = fmt.Errorf("репозиторий баланса не задан")

// BalanceRepositoryAdapter добавляет стратегию повторов для репозитория баланса.
type BalanceRepositoryAdapter struct {
	repo     ports.BalanceRepository
	strategy *retry.Strategy
}

// NewBalanceRepositoryAdapter создаёт адаптер репозитория баланса с поддержкой retry.
func NewBalanceRepositoryAdapter(repo ports.BalanceRepository, strategy *retry.Strategy) (*BalanceRepositoryAdapter, error) {
	if repo == nil {
		return nil, ErrBalanceRepoNil
	}

	return &BalanceRepositoryAdapter{
		repo:     repo,
		strategy: strategy,
	}, nil
}

// GetByUserID возвращает баланс пользователя.
func (a *BalanceRepositoryAdapter) GetByUserID(ctx context.Context, userID int64) (*domain.Balance, error) {
	var balance *domain.Balance
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		balance, err = a.repo.GetByUserID(ctx, userID)
		return err
	})
	return balance, err
}

// CreateForUser создаёт запись баланса для пользователя.
func (a *BalanceRepositoryAdapter) CreateForUser(ctx context.Context, userID int64) error {
	return a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		return a.repo.CreateForUser(ctx, userID)
	})
}

// AddAccrual добавляет начисление к балансу пользователя.
func (a *BalanceRepositoryAdapter) AddAccrual(ctx context.Context, userID int64, amount decimal.Decimal) error {
	return a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		return a.repo.AddAccrual(ctx, userID, amount)
	})
}

// Withdraw выполняет списание средств с баланса пользователя.
func (a *BalanceRepositoryAdapter) Withdraw(ctx context.Context, userID int64, orderNumber string, amount decimal.Decimal) error {
	return a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		return a.repo.Withdraw(ctx, userID, orderNumber, amount)
	})
}
