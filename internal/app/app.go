package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	httpapi "github.com/arvaliullin/gophermart/internal/api/http"
	"github.com/arvaliullin/gophermart/internal/api/http/client/accrual"
	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/config"
	accrualworker "github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/arvaliullin/gophermart/internal/core/services/auth"
	"github.com/arvaliullin/gophermart/internal/core/services/balance"
	"github.com/arvaliullin/gophermart/internal/core/services/order"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	retryadapter "github.com/arvaliullin/gophermart/internal/repository/retry"
	"github.com/rs/zerolog"
)

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
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLoadConfig, err)
	}

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	logger.Info().
		Str("address", cfg.RunAddress).
		Str("database", cfg.DatabaseURI).
		Str("accrual", cfg.AccrualSystemAddress).
		Msg(msgConfigLoaded)

	db, err := postgres.NewDB(ctx, cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectDB, err)
	}

	jwtManager := jwt.NewManager(cfg.JWTSecret)

	retryStrategy := retry.NewStrategy(retry.DefaultDelays, postgres.IsConnectionRetryable)

	userRepo, err := retryadapter.NewUserRepositoryAdapter(
		postgres.NewUserRepository(db.Pool), retryStrategy)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateRetryRepo, err)
	}

	orderRepo, err := retryadapter.NewOrderRepositoryAdapter(
		postgres.NewOrderRepository(db.Pool), retryStrategy)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateRetryRepo, err)
	}

	balanceRepo, err := retryadapter.NewBalanceRepositoryAdapter(
		postgres.NewBalanceRepository(db.Pool), retryStrategy)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateRetryRepo, err)
	}

	withdrawalRepo, err := retryadapter.NewWithdrawalRepositoryAdapter(
		postgres.NewWithdrawalRepository(db.Pool), retryStrategy)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateRetryRepo, err)
	}

	authService := auth.NewService(userRepo, balanceRepo, jwtManager)
	orderService := order.NewService(orderRepo)
	balanceService := balance.NewService(balanceRepo, withdrawalRepo)

	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress)
	worker := accrualworker.NewWorker(orderRepo, balanceRepo, accrualClient, logger)

	authHandler := handlers.NewAuthHandler(authService)
	orderHandler := handlers.NewOrderHandler(orderService)
	balanceHandler := handlers.NewBalanceHandler(balanceService)
	withdrawalHandler := handlers.NewWithdrawalHandler(balanceService)

	router := httpapi.NewRouter(&httpapi.RouterConfig{
		AuthHandler:       authHandler,
		OrderHandler:      orderHandler,
		BalanceHandler:    balanceHandler,
		WithdrawalHandler: withdrawalHandler,
		JWTManager:        jwtManager,
		Logger:            logger,
	})

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	return &App{
		logger:        logger,
		config:        cfg,
		server:        server,
		db:            db,
		accrualWorker: worker,
	}, nil
}

// Logger возвращает логгер приложения.
func (a *App) Logger() *zerolog.Logger {
	return &a.logger
}

// Config возвращает конфигурацию приложения.
func (a *App) Config() *config.AppConfig {
	return a.config
}
