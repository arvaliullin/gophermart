package domain

import "time"

// OrderStatus определяет статус обработки заказа.
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

// Order представляет заказ пользователя в системе лояльности.
type Order struct {
	ID         int64
	UserID     int64
	Number     string
	Status     OrderStatus
	Accrual    *float64
	UploadedAt time.Time
}

// IsFinal возвращает true, если статус заказа является конечным.
func (s OrderStatus) IsFinal() bool {
	return s == OrderStatusInvalid || s == OrderStatusProcessed
}

