package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = fmt.Errorf("недействительный токен")
	ErrExpiredToken = fmt.Errorf("токен истёк")
)

const (
	tokenExpiration = 24 * time.Hour
)

// Manager управляет генерацией и валидацией JWT токенов.
type Manager struct {
	secret []byte
}

// Claims содержит данные JWT токена.
type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

// NewManager создаёт новый менеджер JWT токенов.
func NewManager(secret string) *Manager {
	return &Manager{
		secret: []byte(secret),
	}
}

// GenerateToken создаёт новый JWT токен для пользователя.
func (m *Manager) GenerateToken(userID int64) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken извлекает ID пользователя из JWT токена.
func (m *Manager) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})

	if err != nil {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	return claims.UserID, nil
}
