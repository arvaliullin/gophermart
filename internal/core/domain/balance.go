package domain

// Balance представляет баланс пользователя в системе лояльности.
type Balance struct {
	UserID    int64
	Current   float64
	Withdrawn float64
}

