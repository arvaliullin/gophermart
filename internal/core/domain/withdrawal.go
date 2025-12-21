package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Withdrawal представляет операцию списания баллов.
type Withdrawal struct {
	ID          int64
	UserID      int64
	OrderNumber string
	Sum         decimal.Decimal
	ProcessedAt time.Time
}
