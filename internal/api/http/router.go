package http

import (
	"net/http"

	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// RouterConfig содержит зависимости для настройки роутера.
type RouterConfig struct {
	AuthHandler       *handlers.AuthHandler
	OrderHandler      *handlers.OrderHandler
	BalanceHandler    *handlers.BalanceHandler
	WithdrawalHandler *handlers.WithdrawalHandler
	JWTManager        *jwt.Manager
	Logger            zerolog.Logger
}

// NewRouter создаёт и настраивает HTTP роутер.
func NewRouter(cfg *RouterConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logging(cfg.Logger))
	router.Use(middleware.GzipDecompress())
	router.Use(middleware.GzipCompress())

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Post("/api/user/register", cfg.AuthHandler.Register)
	router.Post("/api/user/login", cfg.AuthHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTManager))

		r.Post("/api/user/orders", cfg.OrderHandler.Submit)
		r.Get("/api/user/orders", cfg.OrderHandler.List)
		r.Get("/api/user/balance", cfg.BalanceHandler.Get)
		r.Post("/api/user/balance/withdraw", cfg.BalanceHandler.Withdraw)
		r.Get("/api/user/withdrawals", cfg.WithdrawalHandler.List)
	})

	return router
}
