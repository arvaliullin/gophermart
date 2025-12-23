package accrual

import "fmt"

var (
	ErrServiceUnavailable = fmt.Errorf("сервис начислений недоступен")
	ErrBuildURL           = fmt.Errorf("ошибка формирования URL")
)
