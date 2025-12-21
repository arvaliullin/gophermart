package order

import (
	"context"
	"errors"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/luhn"
)

// Service реализует бизнес-логику управления заказами.
type Service struct {
	orderRepo ports.OrderRepository
}

// NewService создаёт новый сервис заказов.
func NewService(orderRepo ports.OrderRepository) *Service {
	return &Service{
		orderRepo: orderRepo,
	}
}

// SubmitOrder добавляет новый заказ для пользователя.
// Возвращает true, если заказ уже существовал у этого пользователя.
func (s *Service) SubmitOrder(ctx context.Context, userID int64, number string) (bool, error) {
	if !luhn.IsValid(number) {
		return false, domain.ErrInvalidOrderNumber
	}

	_, err := s.orderRepo.Create(ctx, userID, number)
	if err != nil {
		if errors.Is(err, domain.ErrOrderAlreadyExists) {
			existingOrder, getErr := s.orderRepo.GetByNumber(ctx, number)
			if getErr != nil {
				return false, getErr
			}

			if existingOrder.UserID == userID {
				return true, nil
			}
			return false, domain.ErrOrderBelongsToOther
		}
		return false, err
	}

	return false, nil
}

// GetUserOrders возвращает все заказы пользователя.
func (s *Service) GetUserOrders(ctx context.Context, userID int64) ([]*domain.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID)
}

