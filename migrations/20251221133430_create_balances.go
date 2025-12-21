package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateBalances, downCreateBalances)
}

func upCreateBalances(ctx context.Context, tx *sql.Tx) error {
	query := `
		CREATE TABLE IF NOT EXISTS balances (
			user_id   BIGINT PRIMARY KEY REFERENCES users(id),
			current   DECIMAL(15, 2) NOT NULL DEFAULT 0,
			withdrawn DECIMAL(15, 2) NOT NULL DEFAULT 0
		)
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateBalances(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS balances`)
	return err
}
