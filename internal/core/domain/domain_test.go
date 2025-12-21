package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name     string
		status   OrderStatus
		expected bool
	}{
		{
			name:     "NEW не является конечным",
			status:   OrderStatusNew,
			expected: false,
		},
		{
			name:     "PROCESSING не является конечным",
			status:   OrderStatusProcessing,
			expected: false,
		},
		{
			name:     "INVALID является конечным",
			status:   OrderStatusInvalid,
			expected: true,
		},
		{
			name:     "PROCESSED является конечным",
			status:   OrderStatusProcessed,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsFinal())
		})
	}
}
