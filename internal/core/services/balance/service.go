package balance

import (
	"context"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/luhn"
	"github.com/shopspring/decimal"
)

// Service реализует бизнес-логику управления балансом.
type Service struct {
	balanceRepo    ports.BalanceRepository
	withdrawalRepo ports.WithdrawalRepository
}

// NewService создаёт новый сервис баланса.
func NewService(balanceRepo ports.BalanceRepository, withdrawalRepo ports.WithdrawalRepository) *Service {
	return &Service{
		balanceRepo:    balanceRepo,
		withdrawalRepo: withdrawalRepo,
	}
}

// GetBalance возвращает баланс пользователя.
func (s *Service) GetBalance(ctx context.Context, userID int64) (*domain.Balance, error) {
	return s.balanceRepo.GetByUserID(ctx, userID)
}

// Withdraw выполняет списание средств с баланса пользователя.
func (s *Service) Withdraw(ctx context.Context, userID int64, orderNumber string, amount decimal.Decimal) error {
	if !luhn.IsValid(orderNumber) {
		return domain.ErrInvalidOrderNumber
	}

	return s.balanceRepo.Withdraw(ctx, userID, orderNumber, amount)
}

// GetWithdrawals возвращает все списания пользователя.
func (s *Service) GetWithdrawals(ctx context.Context, userID int64) ([]*domain.Withdrawal, error) {
	return s.withdrawalRepo.GetByUserID(ctx, userID)
}
