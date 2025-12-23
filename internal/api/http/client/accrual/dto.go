package accrual

import "github.com/shopspring/decimal"

type accrualResponse struct {
	Order   string           `json:"order"`
	Status  string           `json:"status"`
	Accrual *decimal.Decimal `json:"accrual,omitempty"`
}
