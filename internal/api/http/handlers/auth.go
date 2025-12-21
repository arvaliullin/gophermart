package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arvaliullin/gophermart/internal/api/http/dto"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
)

var (
	ErrInvalidRequestFormat = fmt.Errorf("неверный формат запроса")
)

const (
	authCookieName = "auth_token"
)

// AuthHandler обрабатывает HTTP запросы аутентификации.
type AuthHandler struct {
	authService ports.AuthService
}

// NewAuthHandler создаёт новый обработчик аутентификации.
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register обрабатывает регистрацию нового пользователя.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

// Login обрабатывает аутентификацию пользователя.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
	w.Header().Set("Authorization", "Bearer "+token)
}
