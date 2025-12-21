package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateWithdrawals, downCreateWithdrawals)
}

func upCreateWithdrawals(ctx context.Context, tx *sql.Tx) error {
	query := `
		CREATE TABLE IF NOT EXISTS withdrawals (
			id           BIGSERIAL PRIMARY KEY,
			user_id      BIGINT NOT NULL REFERENCES users(id),
			order_number VARCHAR(255) NOT NULL,
			sum          DECIMAL(15, 2) NOT NULL,
			processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_withdrawals_user_id ON withdrawals(user_id)
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateWithdrawals(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS withdrawals`)
	return err
}
