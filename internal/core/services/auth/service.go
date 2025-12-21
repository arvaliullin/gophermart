package auth

import (
	"context"
	"errors"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Service реализует бизнес-логику аутентификации пользователей.
type Service struct {
	userRepo    ports.UserRepository
	balanceRepo ports.BalanceRepository
	jwtManager  *jwt.Manager
}

// NewService создаёт новый сервис аутентификации.
func NewService(userRepo ports.UserRepository, balanceRepo ports.BalanceRepository, jwtManager *jwt.Manager) *Service {
	return &Service{
		userRepo:    userRepo,
		balanceRepo: balanceRepo,
		jwtManager:  jwtManager,
	}
}

// Register регистрирует нового пользователя и возвращает JWT токен.
func (s *Service) Register(ctx context.Context, login, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.Create(ctx, login, string(hashedPassword))
	if err != nil {
		return "", err
	}

	if err := s.balanceRepo.CreateForUser(ctx, user.ID); err != nil {
		return "", err
	}

	return s.jwtManager.GenerateToken(user.ID)
}

// Login аутентифицирует пользователя и возвращает JWT токен.
func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.jwtManager.GenerateToken(user.ID)
}
