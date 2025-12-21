//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsConnectionRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name: "connection exception",
			err: &pgconn.PgError{
				Code: pgerrcode.ConnectionException,
			},
			expected: true,
		},
		{
			name: "connection failure",
			err: &pgconn.PgError{
				Code: pgerrcode.ConnectionFailure,
			},
			expected: true,
		},
		{
			name: "non-retryable pg error",
			err: &pgconn.PgError{
				Code: pgerrcode.UniqueViolation,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := postgres.IsConnectionRetryable(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
