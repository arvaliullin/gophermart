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

// UserRepository реализует интерфейс ports.UserRepository для PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository создаёт новый репозиторий пользователей.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create создаёт нового пользователя и возвращает его.
func (r *UserRepository) Create(ctx context.Context, login, passwordHash string) (*domain.User, error) {
	query := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id, login, password, created_at
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

// GetByLogin возвращает пользователя по логину.
func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `
		SELECT id, login, password, created_at
		FROM users
		WHERE login = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByID возвращает пользователя по ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, login, password, created_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
