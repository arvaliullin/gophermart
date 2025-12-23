package app

import "fmt"

// ErrLoadConfig возвращается при ошибке загрузки конфигурации.
var ErrLoadConfig = fmt.Errorf("ошибка загрузки конфигурации")

// ErrConnectDB возвращается при ошибке подключения к базе данных.
var ErrConnectDB = fmt.Errorf("ошибка подключения к БД")

// ErrCreateRetryRepo возвращается при ошибке создания репозитория с retry.
var ErrCreateRetryRepo = fmt.Errorf("ошибка создания репозитория с retry")

const (
	msgConfigLoaded       = "конфигурация загружена"
	msgServerStarting     = "запуск HTTP сервера"
	msgServerError        = "ошибка HTTP сервера"
	msgShuttingDown       = "завершение работы приложения"
	msgServerStopError    = "ошибка остановки HTTP сервера"
	msgDBConnectionClosed = "соединение с БД закрыто"
)

