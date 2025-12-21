//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/arvaliullin/gophermart/internal/repository/postgres/testhelpers"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testContainer *testhelpers.PostgresContainer
	testPool      *pgxpool.Pool
)

// TestMain настраивает тестовый контейнер PostgreSQL для всех интеграционных тестов.
func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = testhelpers.NewPostgresContainer(ctx)
	if err != nil {
		panic("не удалось запустить тестовый контейнер: " + err.Error())
	}

	testPool, err = pgxpool.New(ctx, testContainer.DSN)
	if err != nil {
		_ = testContainer.Terminate(ctx)
		panic("не удалось создать пул соединений: " + err.Error())
	}

	code := m.Run()

	testPool.Close()
	_ = testContainer.Terminate(ctx)

	os.Exit(code)
}

// setupTest очищает таблицы перед каждым тестом.
func setupTest(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	if err := testContainer.CleanupTables(ctx); err != nil {
		t.Fatalf("не удалось очистить таблицы: %v", err)
	}
}
