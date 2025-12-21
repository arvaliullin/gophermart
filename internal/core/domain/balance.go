package domain

import "github.com/shopspring/decimal"

// Balance представляет баланс пользователя в системе лояльности.
type Balance struct {
	UserID    int64
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}
