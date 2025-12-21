package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateOrders, downCreateOrders)
}

func upCreateOrders(ctx context.Context, tx *sql.Tx) error {
	query := `
		CREATE TABLE IF NOT EXISTS orders (
			id          BIGSERIAL PRIMARY KEY,
			user_id     BIGINT NOT NULL REFERENCES users(id),
			number      VARCHAR(255) UNIQUE NOT NULL,
			status      VARCHAR(50) NOT NULL DEFAULT 'NEW',
			accrual     DECIMAL(15, 2),
			uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateOrders(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS orders`)
	return err
}
