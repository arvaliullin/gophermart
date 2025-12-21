package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	loggingMiddleware := middleware.Logging(logger)
	wrappedHandler := loggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test response", rr.Body.String())

	logOutput := buf.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test-path")
	assert.Contains(t, logOutput, "200")
}

func TestLogging_DifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := zerolog.New(&buf)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			loggingMiddleware := middleware.Logging(logger)
			wrappedHandler := loggingMiddleware(handler)

			req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestLogging_ResponseSize(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	responseBody := "hello world!"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(responseBody))
	})

	loggingMiddleware := middleware.Logging(logger)
	wrappedHandler := loggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, responseBody, rr.Body.String())
}
