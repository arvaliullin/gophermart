package dto

import "github.com/arvaliullin/gophermart/internal/core/domain"

// BalanceResponse представляет ответ с информацией о балансе.
type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// WithdrawRequest представляет запрос на списание средств.
type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// IsValid проверяет корректность данных запроса.
func (r *WithdrawRequest) IsValid() bool {
	return r.Order != "" && r.Sum > 0
}

// FromDomainBalance преобразует доменный баланс в DTO.
func FromDomainBalance(balance *domain.Balance) *BalanceResponse {
	return &BalanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}
}

