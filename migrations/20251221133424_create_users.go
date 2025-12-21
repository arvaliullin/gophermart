package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsers, downCreateUsers)
}

func upCreateUsers(ctx context.Context, tx *sql.Tx) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id         BIGSERIAL PRIMARY KEY,
			login      VARCHAR(255) UNIQUE NOT NULL,
			password   VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS users`)
	return err
}
