package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/arvaliullin/gophermart/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer представляет тестовый контейнер PostgreSQL.
type PostgresContainer struct {
	*postgres.PostgresContainer
	DSN string
}

// NewPostgresContainer создаёт новый контейнер PostgreSQL для тестирования.
func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("запуск контейнера postgres: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("получение строки подключения: %w", err)
	}

	if err := runMigrations(ctx, dsn); err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("применение миграций: %w", err)
	}

	return &PostgresContainer{
		PostgresContainer: container,
		DSN:               dsn,
	}, nil
}

// Terminate останавливает и удаляет контейнер.
func (c *PostgresContainer) Terminate(ctx context.Context) error {
	return c.PostgresContainer.Terminate(ctx)
}

// CleanupTables очищает все таблицы для повторного использования контейнера между тестами.
func (c *PostgresContainer) CleanupTables(ctx context.Context) error {
	db, err := sql.Open("pgx", c.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	tables := []string{"withdrawals", "orders", "balances", "users"}
	for _, table := range tables {
		if _, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			return fmt.Errorf("очистка таблицы %s: %w", table, err)
		}
	}

	return nil
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

	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..")
	migrationsDir := filepath.Join(projectRoot, "migrations")

	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		return fmt.Errorf("применение миграций: %w", err)
	}

	return nil
}
