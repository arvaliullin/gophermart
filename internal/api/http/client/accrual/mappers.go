package accrual

import (
	"strconv"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
)

func mapStatus(status string) domain.OrderStatus {
	statusMap := map[string]domain.OrderStatus{
		"REGISTERED": domain.OrderStatusNew,
		"PROCESSING": domain.OrderStatusProcessing,
		"INVALID":    domain.OrderStatusInvalid,
		"PROCESSED":  domain.OrderStatusProcessed,
	}

	if s, ok := statusMap[status]; ok {
		return s
	}
	return domain.OrderStatusNew
}

func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 60 * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 60 * time.Second
	}

	return time.Duration(seconds) * time.Second
}
