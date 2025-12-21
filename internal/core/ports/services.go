package ports

//go:generate mockgen -source=services.go -destination=mocks/services_mock.go -package=mocks

import (
	"context"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/shopspring/decimal"
)

// AuthService определяет контракт сервиса аутентификации.
type AuthService interface {
	Register(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
}

// OrderService определяет контракт сервиса заказов.
type OrderService interface {
	SubmitOrder(ctx context.Context, userID int64, number string) (bool, error)
	GetUserOrders(ctx context.Context, userID int64) ([]*domain.Order, error)
}

// BalanceService определяет контракт сервиса баланса.
type BalanceService interface {
	GetBalance(ctx context.Context, userID int64) (*domain.Balance, error)
	Withdraw(ctx context.Context, userID int64, orderNumber string, amount decimal.Decimal) error
	GetWithdrawals(ctx context.Context, userID int64) ([]*domain.Withdrawal, error)
}
