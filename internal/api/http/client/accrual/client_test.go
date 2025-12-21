package accrual_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/gophermart/internal/api/http/client/accrual"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	accrualservice "github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetOrderAccrual_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/orders/12345678903", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"order":   "12345678903",
			"status":  "PROCESSED",
			"accrual": 500.0,
		})
	}))
	defer server.Close()

	client := accrual.NewClient(server.URL)

	resp, err := client.GetOrderAccrual(context.Background(), "12345678903")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "12345678903", resp.Order)
	assert.Equal(t, domain.OrderStatusProcessed, resp.Status)
	assert.Equal(t, 500.0, resp.Accrual)
}

func TestClient_GetOrderAccrual_Processing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"order":  "12345678903",
			"status": "PROCESSING",
		})
	}))
	defer server.Close()

	client := accrual.NewClient(server.URL)

	resp, err := client.GetOrderAccrual(context.Background(), "12345678903")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, domain.OrderStatusProcessing, resp.Status)
}

func TestClient_GetOrderAccrual_NoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := accrual.NewClient(server.URL)

	resp, err := client.GetOrderAccrual(context.Background(), "12345678903")

	require.NoError(t, err)
	assert.Nil(t, resp)
}

func TestClient_GetOrderAccrual_TooManyRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := accrual.NewClient(server.URL)

	_, err := client.GetOrderAccrual(context.Background(), "12345678903")

	require.Error(t, err)
	var retryErr *accrualservice.RetryAfterError
	assert.ErrorAs(t, err, &retryErr)
	assert.Equal(t, 60, int(retryErr.Duration.Seconds()))
}

func TestClient_GetOrderAccrual_ServiceUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := accrual.NewClient(server.URL)

	_, err := client.GetOrderAccrual(context.Background(), "12345678903")

	require.Error(t, err)
	assert.ErrorIs(t, err, accrual.ErrServiceUnavailable)
}

func TestClient_GetOrderAccrual_StatusMapping(t *testing.T) {
	tests := []struct {
		accrualStatus  string
		expectedStatus domain.OrderStatus
	}{
		{"REGISTERED", domain.OrderStatusNew},
		{"PROCESSING", domain.OrderStatusProcessing},
		{"INVALID", domain.OrderStatusInvalid},
		{"PROCESSED", domain.OrderStatusProcessed},
		{"UNKNOWN", domain.OrderStatusNew},
	}

	for _, tt := range tests {
		t.Run(tt.accrualStatus, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"order":  "12345678903",
					"status": tt.accrualStatus,
				})
			}))
			defer server.Close()

			client := accrual.NewClient(server.URL)

			resp, err := client.GetOrderAccrual(context.Background(), "12345678903")

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.Status)
		})
	}
}
