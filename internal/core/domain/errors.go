package domain

import "fmt"

// Ошибки бизнес-логики системы лояльности.
var (
	// ErrUserNotFound возвращается когда пользователь не найден в системе.
	ErrUserNotFound = fmt.Errorf("пользователь не найден")
	// ErrUserAlreadyExists возвращается при попытке регистрации с занятым логином.
	ErrUserAlreadyExists = fmt.Errorf("пользователь уже существует")
	// ErrInvalidCredentials возвращается при неверной паре логин/пароль.
	ErrInvalidCredentials = fmt.Errorf("неверный логин или пароль")
	// ErrOrderNotFound возвращается когда заказ не найден в системе.
	ErrOrderNotFound = fmt.Errorf("заказ не найден")
	// ErrOrderAlreadyExists возвращается при попытке создать дублирующий заказ.
	ErrOrderAlreadyExists = fmt.Errorf("заказ уже существует")
	// ErrOrderBelongsToOther возвращается когда заказ принадлежит другому пользователю.
	ErrOrderBelongsToOther = fmt.Errorf("заказ принадлежит другому пользователю")
	// ErrInvalidOrderNumber возвращается при невалидном номере заказа.
	ErrInvalidOrderNumber = fmt.Errorf("неверный формат номера заказа")
	// ErrInsufficientBalance возвращается при недостатке средств для списания.
	ErrInsufficientBalance = fmt.Errorf("недостаточно средств на счёте")
	// ErrWithdrawalNotFound возвращается когда операция списания не найдена.
	ErrWithdrawalNotFound = fmt.Errorf("списание не найдено")
)
