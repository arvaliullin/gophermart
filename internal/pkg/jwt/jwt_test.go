package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GenerateAndParseToken(t *testing.T) {
	manager := NewManager("test-secret-key")

	tests := []struct {
		name   string
		userID int64
	}{
		{
			name:   "valid user id",
			userID: 123,
		},
		{
			name:   "zero user id",
			userID: 0,
		},
		{
			name:   "large user id",
			userID: 9999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateToken(tt.userID)
			require.NoError(t, err)
			assert.NotEmpty(t, token)

			parsedUserID, err := manager.ParseToken(token)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, parsedUserID)
		})
	}
}

func TestManager_ParseToken_InvalidToken(t *testing.T) {
	manager := NewManager("test-secret-key")

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid format",
			token: "not-a-valid-token",
		},
		{
			name:  "tampered token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ParseToken(tt.token)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidToken)
		})
	}
}

func TestManager_ParseToken_WrongSecret(t *testing.T) {
	manager1 := NewManager("secret-1")
	manager2 := NewManager("secret-2")

	token, err := manager1.GenerateToken(123)
	require.NoError(t, err)

	_, err = manager2.ParseToken(token)
	assert.Error(t, err)
}
