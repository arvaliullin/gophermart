package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arvaliullin/gophermart/internal/app"
)

const (
	_errorMessage = "failed to run Gophermart App"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	gophermartApp, err := app.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := gophermartApp.Run(ctx); err != nil {
		gophermartApp.Logger().Fatal().Err(err).Msg(_errorMessage)
	}
}
