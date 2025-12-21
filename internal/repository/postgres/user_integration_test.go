//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(testPool)

	t.Run("успешное создание пользователя", func(t *testing.T) {
		user, err := repo.Create(ctx, "testuser", "hashedpassword")
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, "testuser", user.Login)
		assert.Equal(t, "hashedpassword", user.Password)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("ошибка при дублировании логина", func(t *testing.T) {
		_, err := repo.Create(ctx, "duplicate", "pass1")
		require.NoError(t, err)

		_, err = repo.Create(ctx, "duplicate", "pass2")
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})
}

func TestUserRepository_GetByLogin(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(testPool)

	t.Run("успешное получение пользователя по логину", func(t *testing.T) {
		created, err := repo.Create(ctx, "findme", "password")
		require.NoError(t, err)

		found, err := repo.GetByLogin(ctx, "findme")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, created.Login, found.Login)
		assert.Equal(t, created.Password, found.Password)
	})

	t.Run("ошибка при несуществующем логине", func(t *testing.T) {
		_, err := repo.GetByLogin(ctx, "nonexistent")
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	setupTest(t)
	ctx := context.Background()
	repo := postgres.NewUserRepository(testPool)

	t.Run("успешное получение пользователя по ID", func(t *testing.T) {
		created, err := repo.Create(ctx, "byid", "password")
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, created.Login, found.Login)
	})

	t.Run("ошибка при несуществующем ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 999999)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}
