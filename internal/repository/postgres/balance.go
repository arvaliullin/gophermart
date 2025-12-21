package postgres

import (
	"context"
	"errors"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// BalanceRepository реализует интерфейс ports.BalanceRepository для PostgreSQL.
type BalanceRepository struct {
	pool *pgxpool.Pool
}

// NewBalanceRepository создаёт новый репозиторий баланса.
func NewBalanceRepository(pool *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{pool: pool}
}

// GetByUserID возвращает баланс пользователя.
func (r *BalanceRepository) GetByUserID(ctx context.Context, userID int64) (*domain.Balance, error) {
	query := `
		SELECT user_id, current, withdrawn
		FROM balances
		WHERE user_id = $1
	`

	var balance domain.Balance
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&balance.UserID,
		&balance.Current,
		&balance.Withdrawn,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.Balance{UserID: userID}, nil
		}
		return nil, err
	}

	return &balance, nil
}

// CreateForUser создаёт запись баланса для пользователя.
func (r *BalanceRepository) CreateForUser(ctx context.Context, userID int64) error {
	query := `
		INSERT INTO balances (user_id, current, withdrawn)
		VALUES ($1, 0, 0)
		ON CONFLICT (user_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

// AddAccrual добавляет начисление к балансу пользователя.
func (r *BalanceRepository) AddAccrual(ctx context.Context, userID int64, amount decimal.Decimal) error {
	query := `
		INSERT INTO balances (user_id, current, withdrawn)
		VALUES ($1, $2, 0)
		ON CONFLICT (user_id) DO UPDATE
		SET current = balances.current + EXCLUDED.current
	`

	_, err := r.pool.Exec(ctx, query, userID, amount)
	return err
}

// Withdraw выполняет списание средств с баланса пользователя.
func (r *BalanceRepository) Withdraw(ctx context.Context, userID int64, orderNumber string, amount decimal.Decimal) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var currentBalance decimal.Decimal
	err = tx.QueryRow(ctx, `
		SELECT current FROM balances WHERE user_id = $1 FOR UPDATE
	`, userID).Scan(&currentBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInsufficientBalance
		}
		return err
	}

	if currentBalance.LessThan(amount) {
		return domain.ErrInsufficientBalance
	}

	_, err = tx.Exec(ctx, `
		UPDATE balances
		SET current = current - $1, withdrawn = withdrawn + $1
		WHERE user_id = $2
	`, amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO withdrawals (user_id, order_number, sum)
		VALUES ($1, $2, $3)
	`, userID, orderNumber, amount)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
