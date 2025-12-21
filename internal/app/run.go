package app

import (
	"context"
	"net/http"
)

// Run запускает приложение и ожидает сигнала завершения.
func (a *App) Run(ctx context.Context) error {
	go a.accrualWorker.Run(ctx)

	go func() {
		a.logger.Info().
			Str("address", a.server.Addr).
			Msg(msgServerStarting)

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error().
				Err(err).
				Msg(msgServerError)
		}
	}()

	<-ctx.Done()

	a.logger.Info().Msg(msgShuttingDown)

	if err := a.server.Shutdown(context.Background()); err != nil {
		a.logger.Error().
			Err(err).
			Msg(msgServerStopError)
	}

	a.db.Close()
	a.logger.Info().Msg(msgDBConnectionClosed)

	return nil
}
