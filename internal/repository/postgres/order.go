package postgres

import (
	"context"
	"errors"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepository реализует интерфейс ports.OrderRepository для PostgreSQL.
type OrderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository создаёт новый репозиторий заказов.
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// Create создаёт новый заказ.
func (r *OrderRepository) Create(ctx context.Context, userID int64, number string) (*domain.Order, error) {
	query := `
		INSERT INTO orders (user_id, number, status)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, number, status, accrual, uploaded_at
	`

	var order domain.Order
	err := r.pool.QueryRow(ctx, query, userID, number, domain.OrderStatusNew).Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrOrderAlreadyExists
		}
		return nil, err
	}

	return &order, nil
}

// GetByNumber возвращает заказ по номеру.
func (r *OrderRepository) GetByNumber(ctx context.Context, number string) (*domain.Order, error) {
	query := `
		SELECT id, user_id, number, status, accrual, uploaded_at
		FROM orders
		WHERE number = $1
	`

	var order domain.Order
	err := r.pool.QueryRow(ctx, query, number).Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

// GetByUserID возвращает все заказы пользователя, отсортированные по дате загрузки.
func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Order, error) {
	query := `
		SELECT id, user_id, number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// GetPendingOrders возвращает заказы со статусами NEW или PROCESSING.
func (r *OrderRepository) GetPendingOrders(ctx context.Context) ([]*domain.Order, error) {
	query := `
		SELECT id, user_id, number, status, accrual, uploaded_at
		FROM orders
		WHERE status IN ($1, $2)
		ORDER BY uploaded_at ASC
	`

	rows, err := r.pool.Query(ctx, query, domain.OrderStatusNew, domain.OrderStatusProcessing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// UpdateStatus обновляет статус и начисление заказа.
func (r *OrderRepository) UpdateStatus(ctx context.Context, number string, status domain.OrderStatus, accrual *float64) error {
	query := `
		UPDATE orders
		SET status = $1, accrual = $2
		WHERE number = $3
	`

	result, err := r.pool.Exec(ctx, query, status, accrual, number)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

