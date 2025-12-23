package accrual

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/shopspring/decimal"
)

// GetOrderAccrual запрашивает информацию о начислениях по заказу.
func (c *Client) GetOrderAccrual(ctx context.Context, orderNumber string) (*ports.AccrualResponse, error) {
	requestURL, err := url.JoinPath(c.baseURL, "api", "orders", orderNumber)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBuildURL, err)
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&accrualResponse{}).
		Get(requestURL)

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		result := resp.Result().(*accrualResponse)
		accrual := decimal.Zero
		if result.Accrual != nil {
			accrual = *result.Accrual
		}
		return &ports.AccrualResponse{
			Order:   result.Order,
			Status:  mapStatus(result.Status),
			Accrual: accrual,
		}, nil

	case http.StatusNoContent:
		return nil, nil

	case http.StatusTooManyRequests:
		retryAfter := resp.Header().Get("Retry-After")
		duration := parseRetryAfter(retryAfter)
		return nil, &accrual.RetryAfterError{Duration: duration}

	default:
		return nil, ErrServiceUnavailable
	}
}
