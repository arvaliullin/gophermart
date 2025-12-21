package accrual

import (
	"context"
	"sync"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

const (
	defaultPollInterval = 1 * time.Second

	msgWorkerStopping      = "остановка воркера начислений"
	msgGetOrdersError      = "ошибка получения заказов"
	msgRateLimitExceeded   = "превышен лимит запросов к системе начислений"
	msgAccrualRequestError = "ошибка запроса к системе начислений"
	msgUpdateStatusError   = "ошибка обновления статуса заказа"
	msgAccrualError        = "ошибка начисления баллов"
	msgAccrualSuccess      = "баллы успешно начислены"
)

// Worker опрашивает систему начислений и обновляет статусы заказов.
type Worker struct {
	orderRepo     ports.OrderRepository
	balanceRepo   ports.BalanceRepository
	accrualClient ports.AccrualClient
	logger        zerolog.Logger
	pollInterval  time.Duration
	retryAfter    time.Duration
	mu            sync.Mutex
}

// NewWorker создаёт новый воркер опроса системы начислений.
func NewWorker(
	orderRepo ports.OrderRepository,
	balanceRepo ports.BalanceRepository,
	accrualClient ports.AccrualClient,
	logger zerolog.Logger,
) *Worker {
	return &Worker{
		orderRepo:     orderRepo,
		balanceRepo:   balanceRepo,
		accrualClient: accrualClient,
		logger:        logger,
		pollInterval:  defaultPollInterval,
		retryAfter:    0,
	}
}

// Run запускает воркер опроса системы начислений.
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg(msgWorkerStopping)
			return
		case <-ticker.C:
			w.processOrders(ctx)
		}
	}
}

func (w *Worker) processOrders(ctx context.Context) {
	w.mu.Lock()
	if w.retryAfter > 0 {
		w.retryAfter -= w.pollInterval
		w.mu.Unlock()
		return
	}
	w.mu.Unlock()

	orders, err := w.orderRepo.GetPendingOrders(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg(msgGetOrdersError)
		return
	}

	for _, order := range orders {
		if ctx.Err() != nil {
			return
		}

		w.processOrder(ctx, order)
	}
}

func (w *Worker) processOrder(ctx context.Context, order *domain.Order) {
	resp, err := w.accrualClient.GetOrderAccrual(ctx, order.Number)
	if err != nil {
		if retryErr, ok := err.(*RetryAfterError); ok {
			w.mu.Lock()
			w.retryAfter = retryErr.Duration
			w.mu.Unlock()
			w.logger.Warn().
				Dur("retry_after", retryErr.Duration).
				Msg(msgRateLimitExceeded)
			return
		}

		w.logger.Error().
			Err(err).
			Str("order", order.Number).
			Msg(msgAccrualRequestError)
		return
	}

	if resp == nil {
		return
	}

	var accrual *decimal.Decimal
	if resp.Accrual.IsPositive() {
		accrual = &resp.Accrual
	}

	err = w.orderRepo.UpdateStatus(ctx, order.Number, resp.Status, accrual)
	if err != nil {
		w.logger.Error().
			Err(err).
			Str("order", order.Number).
			Msg(msgUpdateStatusError)
		return
	}

	if resp.Status == domain.OrderStatusProcessed && resp.Accrual.IsPositive() {
		err = w.balanceRepo.AddAccrual(ctx, order.UserID, resp.Accrual)
		if err != nil {
			w.logger.Error().
				Err(err).
				Str("order", order.Number).
				Str("accrual", resp.Accrual.String()).
				Msg(msgAccrualError)
			return
		}

		w.logger.Info().
			Str("order", order.Number).
			Str("accrual", resp.Accrual.String()).
			Msg(msgAccrualSuccess)
	}
}
