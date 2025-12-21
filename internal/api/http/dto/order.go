package dto

import (
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
)

// OrderResponse представляет ответ с информацией о заказе.
type OrderResponse struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

// FromDomainOrder преобразует доменный заказ в DTO.
func FromDomainOrder(order *domain.Order) *OrderResponse {
	return &OrderResponse{
		Number:     order.Number,
		Status:     string(order.Status),
		Accrual:    order.Accrual,
		UploadedAt: order.UploadedAt.Format(time.RFC3339),
	}
}

// FromDomainOrders преобразует список доменных заказов в список DTO.
func FromDomainOrders(orders []*domain.Order) []*OrderResponse {
	result := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		result[i] = FromDomainOrder(order)
	}
	return result
}
