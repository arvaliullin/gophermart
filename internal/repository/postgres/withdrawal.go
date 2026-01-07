package postgres

import (
	"context"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
)

type WithdrawalRepository struct {
	client ports.PostgresClient
}

// NewWithdrawalRepository создаёт новый репозиторий списаний.
func NewWithdrawalRepository(client ports.PostgresClient) *WithdrawalRepository {
	return &WithdrawalRepository{client: client}
}

// GetByUserID возвращает все списания пользователя, отсортированные по дате.
func (r *WithdrawalRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Withdrawal, error) {
	query := `
		SELECT id, user_id, order_number, sum, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`

	rows, err := r.client.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []*domain.Withdrawal
	for rows.Next() {
		var w domain.Withdrawal
		err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OrderNumber,
			&w.Sum,
			&w.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, &w)
	}

	return withdrawals, rows.Err()
}
