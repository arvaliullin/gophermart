package dto

import (
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequest_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		request  AuthRequest
		expected bool
	}{
		{
			name:     "валидный запрос",
			request:  AuthRequest{Login: "user", Password: "pass"},
			expected: true,
		},
		{
			name:     "пустой логин",
			request:  AuthRequest{Login: "", Password: "pass"},
			expected: false,
		},
		{
			name:     "пустой пароль",
			request:  AuthRequest{Login: "user", Password: ""},
			expected: false,
		},
		{
			name:     "пустой запрос",
			request:  AuthRequest{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.request.IsValid())
		})
	}
}

func TestWithdrawRequest_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		request  WithdrawRequest
		expected bool
	}{
		{
			name:     "валидный запрос",
			request:  WithdrawRequest{Order: "123", Sum: 100.0},
			expected: true,
		},
		{
			name:     "пустой номер заказа",
			request:  WithdrawRequest{Order: "", Sum: 100.0},
			expected: false,
		},
		{
			name:     "нулевая сумма",
			request:  WithdrawRequest{Order: "123", Sum: 0},
			expected: false,
		},
		{
			name:     "отрицательная сумма",
			request:  WithdrawRequest{Order: "123", Sum: -10},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.request.IsValid())
		})
	}
}

func TestFromDomainBalance(t *testing.T) {
	balance := &domain.Balance{
		UserID:    1,
		Current:   500.50,
		Withdrawn: 100.25,
	}

	response := FromDomainBalance(balance)

	assert.Equal(t, 500.50, response.Current)
	assert.Equal(t, 100.25, response.Withdrawn)
}

func TestFromDomainOrder(t *testing.T) {
	uploadedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	accrual := 250.0
	order := &domain.Order{
		ID:         1,
		UserID:     1,
		Number:     "12345678903",
		Status:     domain.OrderStatusProcessed,
		Accrual:    &accrual,
		UploadedAt: uploadedAt,
	}

	response := FromDomainOrder(order)

	assert.Equal(t, "12345678903", response.Number)
	assert.Equal(t, "PROCESSED", response.Status)
	assert.Equal(t, &accrual, response.Accrual)
	assert.Equal(t, "2024-01-15T10:30:00Z", response.UploadedAt)
}

func TestFromDomainOrder_WithoutAccrual(t *testing.T) {
	uploadedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	order := &domain.Order{
		ID:         1,
		UserID:     1,
		Number:     "12345678903",
		Status:     domain.OrderStatusNew,
		Accrual:    nil,
		UploadedAt: uploadedAt,
	}

	response := FromDomainOrder(order)

	assert.Equal(t, "12345678903", response.Number)
	assert.Equal(t, "NEW", response.Status)
	assert.Nil(t, response.Accrual)
}

func TestFromDomainOrders(t *testing.T) {
	uploadedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	orders := []*domain.Order{
		{Number: "111", Status: domain.OrderStatusNew, UploadedAt: uploadedAt},
		{Number: "222", Status: domain.OrderStatusProcessing, UploadedAt: uploadedAt},
	}

	responses := FromDomainOrders(orders)

	assert.Len(t, responses, 2)
	assert.Equal(t, "111", responses[0].Number)
	assert.Equal(t, "NEW", responses[0].Status)
	assert.Equal(t, "222", responses[1].Number)
	assert.Equal(t, "PROCESSING", responses[1].Status)
}

func TestFromDomainOrders_Empty(t *testing.T) {
	responses := FromDomainOrders([]*domain.Order{})
	assert.Empty(t, responses)
}

func TestFromDomainWithdrawal(t *testing.T) {
	processedAt := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	withdrawal := &domain.Withdrawal{
		ID:          1,
		UserID:      1,
		OrderNumber: "2377225624",
		Sum:         500.0,
		ProcessedAt: processedAt,
	}

	response := FromDomainWithdrawal(withdrawal)

	assert.Equal(t, "2377225624", response.Order)
	assert.Equal(t, 500.0, response.Sum)
	assert.Equal(t, "2024-01-15T12:00:00Z", response.ProcessedAt)
}

func TestFromDomainWithdrawals(t *testing.T) {
	processedAt := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	withdrawals := []*domain.Withdrawal{
		{OrderNumber: "111", Sum: 100.0, ProcessedAt: processedAt},
		{OrderNumber: "222", Sum: 200.0, ProcessedAt: processedAt},
	}

	responses := FromDomainWithdrawals(withdrawals)

	assert.Len(t, responses, 2)
	assert.Equal(t, "111", responses[0].Order)
	assert.Equal(t, 100.0, responses[0].Sum)
	assert.Equal(t, "222", responses[1].Order)
	assert.Equal(t, 200.0, responses[1].Sum)
}

func TestFromDomainWithdrawals_Empty(t *testing.T) {
	responses := FromDomainWithdrawals([]*domain.Withdrawal{})
	assert.Empty(t, responses)
}
