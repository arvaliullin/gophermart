package dto

import (
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/shopspring/decimal"
)

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

// GetSumAsDecimal возвращает сумму как decimal.Decimal.
func (r *WithdrawRequest) GetSumAsDecimal() decimal.Decimal {
	return decimal.NewFromFloat(r.Sum)
}

// FromDomainBalance преобразует доменный баланс в DTO.
func FromDomainBalance(balance *domain.Balance) *BalanceResponse {
	current, _ := balance.Current.Float64()
	withdrawn, _ := balance.Withdrawn.Float64()
	return &BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	}
}
