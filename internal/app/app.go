package app

import (
	"context"
	"fmt"
	"os"

	"github.com/arvaliullin/gophermart/internal/config"
	"github.com/rs/zerolog"
)

// App представляет основное приложение Gophermart.
type App struct {
	logger zerolog.Logger
	config *config.AppConfig
}

func New(ctx context.Context) (*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	return &App{
		logger: logger,
		config: cfg,
	}, nil
}

func (a *App) Logger() *zerolog.Logger {
	return &a.logger
}

func (a *App) Config() *config.AppConfig {
	return a.config
}
