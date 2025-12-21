package domain

import "fmt"

var (
	ErrUserNotFound         = fmt.Errorf("пользователь не найден")
	ErrUserAlreadyExists    = fmt.Errorf("пользователь уже существует")
	ErrInvalidCredentials   = fmt.Errorf("неверный логин или пароль")
	ErrOrderNotFound        = fmt.Errorf("заказ не найден")
	ErrOrderAlreadyExists   = fmt.Errorf("заказ уже существует")
	ErrOrderBelongsToOther  = fmt.Errorf("заказ принадлежит другому пользователю")
	ErrInvalidOrderNumber   = fmt.Errorf("неверный формат номера заказа")
	ErrInsufficientBalance  = fmt.Errorf("недостаточно средств на счёте")
	ErrWithdrawalNotFound   = fmt.Errorf("списание не найдено")
)

