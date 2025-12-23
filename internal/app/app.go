package app

import (
	"context"
	"net/http"

	"github.com/arvaliullin/gophermart/internal/config"
	accrualworker "github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	"github.com/rs/zerolog"
)

// App представляет основное приложение Gophermart.
type App struct {
	logger        zerolog.Logger
	config        *config.AppConfig
	server        *http.Server
	db            *postgres.DB
	accrualWorker *accrualworker.Worker
}

// New создаёт новый экземпляр приложения с инициализированными зависимостями.
func New(ctx context.Context) (*App, error) {
	return NewBuilder(ctx).
		WithConfig().
		WithLogger().
		WithDatabase().
		WithInfrastructure().
		WithRepositories().
		WithServices().
		WithAccrualWorker().
		WithHTTPServer().
		Build()
}

// Logger возвращает логгер приложения.
func (a *App) Logger() *zerolog.Logger {
	return &a.logger
}

// Config возвращает конфигурацию приложения.
func (a *App) Config() *config.AppConfig {
	return a.config
}
