package app

import (
	"context"
)

const (
	_message = "starting app"
)

// Run запускает приложение.
func (a *App) Run(ctx context.Context) error {

	a.logger.Info().
		Msg(_message)

	return nil
}
