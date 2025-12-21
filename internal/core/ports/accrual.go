package ports

import (
	"context"

	"github.com/arvaliullin/gophermart/internal/core/domain"
)

//go:generate mockgen -source=accrual.go -destination=mocks/accrual_mock.go -package=mocks

// AccrualResponse представляет ответ от системы начислений.
type AccrualResponse struct {
	Order   string
	Status  domain.OrderStatus
	Accrual float64
}

// AccrualClient определяет контракт для взаимодействия с системой начислений.
type AccrualClient interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (*AccrualResponse, error)
}
