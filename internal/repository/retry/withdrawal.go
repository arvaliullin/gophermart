package retry

import (
	"context"
	"fmt"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
)

var (
	ErrWithdrawalRepoNil = fmt.Errorf("репозиторий списаний не задан")
)

// WithdrawalRepositoryAdapter добавляет стратегию повторов для репозитория списаний.
type WithdrawalRepositoryAdapter struct {
	repo     ports.WithdrawalRepository
	strategy *retry.Strategy
}

// NewWithdrawalRepositoryAdapter создаёт адаптер репозитория списаний с поддержкой retry.
func NewWithdrawalRepositoryAdapter(repo ports.WithdrawalRepository, strategy *retry.Strategy) (*WithdrawalRepositoryAdapter, error) {
	if repo == nil {
		return nil, ErrWithdrawalRepoNil
	}

	return &WithdrawalRepositoryAdapter{
		repo:     repo,
		strategy: strategy,
	}, nil
}

// GetByUserID возвращает все списания пользователя.
func (a *WithdrawalRepositoryAdapter) GetByUserID(ctx context.Context, userID int64) ([]*domain.Withdrawal, error) {
	var withdrawals []*domain.Withdrawal
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		withdrawals, err = a.repo.GetByUserID(ctx, userID)
		return err
	})
	return withdrawals, err
}

