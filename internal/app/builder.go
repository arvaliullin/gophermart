package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	httpapi "github.com/arvaliullin/gophermart/internal/api/http"
	"github.com/arvaliullin/gophermart/internal/api/http/client/accrual"
	"github.com/arvaliullin/gophermart/internal/api/http/handlers"
	"github.com/arvaliullin/gophermart/internal/config"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	accrualworker "github.com/arvaliullin/gophermart/internal/core/services/accrual"
	"github.com/arvaliullin/gophermart/internal/core/services/auth"
	"github.com/arvaliullin/gophermart/internal/core/services/balance"
	"github.com/arvaliullin/gophermart/internal/core/services/order"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"github.com/arvaliullin/gophermart/internal/pkg/retry"
	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	retryadapter "github.com/arvaliullin/gophermart/internal/repository/retry"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Builder используется для пошагового построения приложения.
type Builder struct {
	ctx    context.Context
	config *config.AppConfig
	logger zerolog.Logger
	db     *postgres.DB

	jwtManager    *jwt.Manager
	retryStrategy *retry.Strategy

	userRepo       ports.UserRepository
	orderRepo      ports.OrderRepository
	balanceRepo    ports.BalanceRepository
	withdrawalRepo ports.WithdrawalRepository

	authService    *auth.Service
	orderService   *order.Service
	balanceService *balance.Service

	accrualClient *accrual.Client
	accrualWorker *accrualworker.Worker

	server *http.Server
}

// NewBuilder создаёт новый экземпляр Builder.
func NewBuilder(ctx context.Context) *Builder {
	return &Builder{
		ctx: ctx,
	}
}

// WithConfig загружает конфигурацию приложения.
func (b *Builder) WithConfig() *Builder {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("%w: %w", ErrLoadConfig, err))
	}
	b.config = cfg
	return b
}

// WithLogger инициализирует логгер.
func (b *Builder) WithLogger() *Builder {
	b.logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	b.logger.Info().
		Str("address", b.config.RunAddress).
		Str("database", b.config.DatabaseURI).
		Str("accrual", b.config.AccrualSystemAddress).
		Msg(msgConfigLoaded)

	return b
}

// WithDatabase устанавливает подключение к базе данных.
func (b *Builder) WithDatabase() *Builder {
	db, err := postgres.NewDB(b.ctx, b.config.DatabaseURI)
	if err != nil {
		panic(fmt.Errorf("%w: %w", ErrConnectDB, err))
	}
	b.db = db
	return b
}

// WithInfrastructure инициализирует инфраструктурные компоненты (JWT, retry).
func (b *Builder) WithInfrastructure() *Builder {
	b.jwtManager = jwt.NewManager(b.config.JWTSecret)
	b.retryStrategy = retry.NewStrategy(retry.DefaultDelays, postgres.IsConnectionRetryable)
	return b
}

// WithRepositories создаёт репозитории с retry адаптерами.
func (b *Builder) WithRepositories() *Builder {
	var err error

	b.userRepo, err = retryadapter.NewUserRepositoryAdapter(
		postgres.NewUserRepository(b.db.Pool), b.retryStrategy)
	if err != nil {
		panic(fmt.Errorf("%w: %w", ErrCreateRetryRepo, err))
	}

	b.orderRepo, err = retryadapter.NewOrderRepositoryAdapter(
		postgres.NewOrderRepository(b.db.Pool), b.retryStrategy)
	if err != nil {
		panic(fmt.Errorf("%w: %w", ErrCreateRetryRepo, err))
	}

	b.balanceRepo, err = retryadapter.NewBalanceRepositoryAdapter(
		postgres.NewBalanceRepository(b.db.Pool), b.retryStrategy)
	if err != nil {
		panic(fmt.Errorf("%w: %w", ErrCreateRetryRepo, err))
	}

	b.withdrawalRepo, err = retryadapter.NewWithdrawalRepositoryAdapter(
		postgres.NewWithdrawalRepository(b.db.Pool), b.retryStrategy)
	if err != nil {
		panic(fmt.Errorf("%w: %w", ErrCreateRetryRepo, err))
	}

	return b
}

// WithServices создаёт бизнес-сервисы.
func (b *Builder) WithServices() *Builder {
	b.authService = auth.NewService(b.userRepo, b.balanceRepo, b.jwtManager)
	b.orderService = order.NewService(b.orderRepo)
	b.balanceService = balance.NewService(b.balanceRepo, b.withdrawalRepo)
	return b
}

// WithAccrualWorker создаёт клиент для accrual системы и воркер.
func (b *Builder) WithAccrualWorker() *Builder {
	httpClient := resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(0)

	b.accrualClient = accrual.NewClient(b.config.AccrualSystemAddress,
		accrual.WithHTTPClient(httpClient))

	b.accrualWorker = accrualworker.NewWorker(
		b.orderRepo,
		b.balanceRepo,
		b.accrualClient,
		b.logger,
	)

	return b
}

// WithHTTPServer создаёт HTTP сервер с роутером.
func (b *Builder) WithHTTPServer() *Builder {
	authHandler := handlers.NewAuthHandler(b.authService)
	orderHandler := handlers.NewOrderHandler(b.orderService)
	balanceHandler := handlers.NewBalanceHandler(b.balanceService)
	withdrawalHandler := handlers.NewWithdrawalHandler(b.balanceService)

	router := httpapi.NewRouter(&httpapi.RouterConfig{
		AuthHandler:       authHandler,
		OrderHandler:      orderHandler,
		BalanceHandler:    balanceHandler,
		WithdrawalHandler: withdrawalHandler,
		JWTManager:        b.jwtManager,
		Logger:            b.logger,
	})

	b.server = &http.Server{
		Addr:    b.config.RunAddress,
		Handler: router,
	}

	return b
}

// Build собирает и возвращает готовый экземпляр приложения.
func (b *Builder) Build() (*App, error) {
	return &App{
		logger:        b.logger,
		config:        b.config,
		server:        b.server,
		db:            b.db,
		accrualWorker: b.accrualWorker,
	}, nil
}
