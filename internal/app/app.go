package app

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// App представляет основное приложение Gophermart.
type App struct {
	logger zerolog.Logger
}

func New(ctx context.Context) *App {
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)
	return &App{logger: logger}
}

func (a *App) Logger() *zerolog.Logger {
	return &a.logger
}
