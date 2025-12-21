package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports/mocks"
	"github.com/arvaliullin/gophermart/internal/core/services/auth"
	"github.com/arvaliullin/gophermart/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	jwtManager := jwt.NewManager("test-secret")

	service := auth.NewService(userRepo, balanceRepo, jwtManager)

	user := &domain.User{
		ID:        1,
		Login:     "testuser",
		CreatedAt: time.Now(),
	}

	userRepo.EXPECT().
		Create(gomock.Any(), "testuser", gomock.Any()).
		Return(user, nil)

	balanceRepo.EXPECT().
		CreateForUser(gomock.Any(), int64(1)).
		Return(nil)

	token, err := service.Register(context.Background(), "testuser", "password123")

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	userID, err := jwtManager.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(1), userID)
}

func TestService_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	jwtManager := jwt.NewManager("test-secret")

	service := auth.NewService(userRepo, balanceRepo, jwtManager)

	userRepo.EXPECT().
		Create(gomock.Any(), "existinguser", gomock.Any()).
		Return(nil, domain.ErrUserAlreadyExists)

	_, err := service.Register(context.Background(), "existinguser", "password123")

	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
}

func TestService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	jwtManager := jwt.NewManager("test-secret")

	service := auth.NewService(userRepo, balanceRepo, jwtManager)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Login:     "testuser",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	token, err := service.Login(context.Background(), "testuser", "password123")

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	userID, err := jwtManager.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(1), userID)
}

func TestService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	jwtManager := jwt.NewManager("test-secret")

	service := auth.NewService(userRepo, balanceRepo, jwtManager)

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "unknown").
		Return(nil, domain.ErrUserNotFound)

	_, err := service.Login(context.Background(), "unknown", "password123")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestService_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	balanceRepo := mocks.NewMockBalanceRepository(ctrl)
	jwtManager := jwt.NewManager("test-secret")

	service := auth.NewService(userRepo, balanceRepo, jwtManager)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Login:     "testuser",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	_, err := service.Login(context.Background(), "testuser", "wrong-password")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}
