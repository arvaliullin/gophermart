package app

import (
	"context"
	"net/http"
	"time"
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error().
			Err(err).
			Msg(msgServerStopError)
	}

	a.db.Close()
	a.logger.Info().Msg(msgDBConnectionClosed)

	return nil
}
