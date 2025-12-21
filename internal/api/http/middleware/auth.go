package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
)

type contextKey string

const userIDKey contextKey = "user_id"
const authCookieName = "auth_token"

// UserIDKey экспортируемый ключ для тестирования.
var UserIDKey = userIDKey

// Auth создаёт middleware для проверки JWT авторизации.
func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
				return
			}

			userID, err := jwtManager.ParseToken(token)
			if err != nil {
				http.Error(w, "недействительный токен", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID извлекает ID пользователя из контекста запроса.
func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	cookie, err := r.Cookie(authCookieName)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}
