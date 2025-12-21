package retry

import (
	"context"
	"fmt"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
)

// ErrUserRepoNil возвращается при попытке создать адаптер с nil репозиторием.
var ErrUserRepoNil = fmt.Errorf("репозиторий пользователей не задан")

// UserRepositoryAdapter добавляет стратегию повторов для репозитория пользователей.
type UserRepositoryAdapter struct {
	repo     ports.UserRepository
	strategy *retry.Strategy
}

// NewUserRepositoryAdapter создаёт адаптер репозитория пользователей с поддержкой retry.
func NewUserRepositoryAdapter(repo ports.UserRepository, strategy *retry.Strategy) (*UserRepositoryAdapter, error) {
	if repo == nil {
		return nil, ErrUserRepoNil
	}

	return &UserRepositoryAdapter{
		repo:     repo,
		strategy: strategy,
	}, nil
}

// Create создаёт нового пользователя.
func (a *UserRepositoryAdapter) Create(ctx context.Context, login, passwordHash string) (*domain.User, error) {
	var user *domain.User
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		user, err = a.repo.Create(ctx, login, passwordHash)
		return err
	})
	return user, err
}

// GetByLogin возвращает пользователя по логину.
func (a *UserRepositoryAdapter) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user *domain.User
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		user, err = a.repo.GetByLogin(ctx, login)
		return err
	})
	return user, err
}

// GetByID возвращает пользователя по ID.
func (a *UserRepositoryAdapter) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user *domain.User
	err := a.strategy.DoWithRetry(ctx, func(ctx context.Context) error {
		var err error
		user, err = a.repo.GetByID(ctx, id)
		return err
	})
	return user, err
}
