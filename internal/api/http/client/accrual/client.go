package accrual

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// HTTPClient определяет интерфейс HTTP клиента для взаимодействия с системой начислений.
type HTTPClient interface {
	R() *resty.Request
}

// Client реализует HTTP клиент для взаимодействия с системой начислений.
type Client struct {
	client  HTTPClient
	baseURL string
}

// ClientOption определяет функциональную опцию для настройки клиента.
type ClientOption func(*clientConfig)

type clientConfig struct {
	httpClient HTTPClient
}

// WithHTTPClient устанавливает пользовательский HTTP клиент.
func WithHTTPClient(client HTTPClient) ClientOption {
	return func(cfg *clientConfig) {
		cfg.httpClient = client
	}
}

// NewClient создаёт новый HTTP клиент системы начислений с указанными опциями.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	cfg := &clientConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.httpClient == nil {
		cfg.httpClient = resty.New().
			SetTimeout(10 * time.Second).
			SetRetryCount(0)
	}

	return &Client{
		client:  cfg.httpClient,
		baseURL: baseURL,
	}
}
