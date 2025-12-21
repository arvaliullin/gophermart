package ports

import (
	"context"

	"github.com/arvaliullin/gophermart/internal/core/domain"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

// UserRepository определяет контракт для работы с пользователями.
type UserRepository interface {
	Create(ctx context.Context, login, passwordHash string) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

// OrderRepository определяет контракт для работы с заказами.
type OrderRepository interface {
	Create(ctx context.Context, userID int64, number string) (*domain.Order, error)
	GetByNumber(ctx context.Context, number string) (*domain.Order, error)
	GetByUserID(ctx context.Context, userID int64) ([]*domain.Order, error)
	GetPendingOrders(ctx context.Context) ([]*domain.Order, error)
	UpdateStatus(ctx context.Context, number string, status domain.OrderStatus, accrual *float64) error
}

// BalanceRepository определяет контракт для работы с балансом.
type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*domain.Balance, error)
	CreateForUser(ctx context.Context, userID int64) error
	AddAccrual(ctx context.Context, userID int64, amount float64) error
	Withdraw(ctx context.Context, userID int64, orderNumber string, amount float64) error
}

// WithdrawalRepository определяет контракт для работы со списаниями.
type WithdrawalRepository interface {
	GetByUserID(ctx context.Context, userID int64) ([]*domain.Withdrawal, error)
}
