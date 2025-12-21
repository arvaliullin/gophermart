package accrual

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

var (
	ErrServiceUnavailable = fmt.Errorf("сервис начислений недоступен")
	ErrBuildURL           = fmt.Errorf("ошибка формирования URL")
)

type accrualResponse struct {
	Order   string          `json:"order"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual,omitempty"`
}

// Client реализует HTTP клиент для взаимодействия с системой начислений.
type Client struct {
	client  *resty.Client
	baseURL string
}

// NewClient создаёт новый HTTP клиент системы начислений.
func NewClient(baseURL string) *Client {
	client := resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(0)

	return &Client{
		client:  client,
		baseURL: baseURL,
	}
}

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
		return &ports.AccrualResponse{
			Order:   result.Order,
			Status:  mapStatus(result.Status),
			Accrual: result.Accrual,
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
