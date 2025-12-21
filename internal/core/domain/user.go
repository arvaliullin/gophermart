package domain

import "time"

// User представляет пользователя системы лояльности.
type User struct {
	ID        int64
	Login     string
	Password  string
	CreatedAt time.Time
}
