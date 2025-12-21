package dto

import (
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
)

// WithdrawalResponse представляет ответ с информацией о списании.
type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

// FromDomainWithdrawal преобразует доменное списание в DTO.
func FromDomainWithdrawal(w *domain.Withdrawal) *WithdrawalResponse {
	return &WithdrawalResponse{
		Order:       w.OrderNumber,
		Sum:         w.Sum,
		ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
	}
}

// FromDomainWithdrawals преобразует список доменных списаний в список DTO.
func FromDomainWithdrawals(withdrawals []*domain.Withdrawal) []*WithdrawalResponse {
	result := make([]*WithdrawalResponse, len(withdrawals))
	for i, w := range withdrawals {
		result[i] = FromDomainWithdrawal(w)
	}
	return result
}

