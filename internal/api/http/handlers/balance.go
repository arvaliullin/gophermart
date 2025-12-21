package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arvaliullin/gophermart/internal/api/http/dto"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
)

// BalanceHandler обрабатывает HTTP запросы управления балансом.
type BalanceHandler struct {
	balanceService ports.BalanceService
}

// NewBalanceHandler создаёт новый обработчик баланса.
func NewBalanceHandler(balanceService ports.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

// Get возвращает текущий баланс пользователя.
func (h *BalanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	bal, err := h.balanceService.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.FromDomainBalance(bal))
}

// Withdraw обрабатывает запрос на списание средств.
func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var req dto.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	err := h.balanceService.Withdraw(r.Context(), userID, req.Order, req.GetSumAsDecimal())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidOrderNumber) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, domain.ErrInsufficientBalance) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
