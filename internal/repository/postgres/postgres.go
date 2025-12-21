package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/arvaliullin/gophermart/migrations"
	shopspring "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var connectionErrorCodes = map[string]struct{}{
	pgerrcode.ConnectionException:                           {},
	pgerrcode.ConnectionDoesNotExist:                        {},
	pgerrcode.ConnectionFailure:                             {},
	pgerrcode.SQLClientUnableToEstablishSQLConnection:       {},
	pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection: {},
	pgerrcode.TransactionResolutionUnknown:                  {},
	pgerrcode.ProtocolViolation:                             {},
}

// DB представляет пул соединений с PostgreSQL.
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB создаёт новое подключение к PostgreSQL с применением миграций.
func NewDB(ctx context.Context, dsn string) (*DB, error) {
	if err := runMigrations(ctx, dsn); err != nil {
		return nil, fmt.Errorf("применение миграций: %w", err)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("парсинг конфигурации: %w", err)
	}

	// Регистрируем поддержку decimal.Decimal для PostgreSQL NUMERIC/DECIMAL
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		shopspring.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("создание пула соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("проверка соединения: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close закрывает пул соединений.
func (db *DB) Close() {
	db.Pool.Close()
}

func runMigrations(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("открытие БД: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("установка диалекта: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("применение миграций: %w", err)
	}

	return nil
}

// IsConnectionRetryable определяет, относится ли ошибка к классу сбоев соединения PostgreSQL.
func IsConnectionRetryable(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		_, ok := connectionErrorCodes[pgErr.Code]
		return ok
	}

	var connectErr *pgconn.ConnectError
	return errors.As(err, &connectErr)
}
