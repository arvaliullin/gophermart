package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/arvaliullin/gophermart/internal/api/http/dto"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/core/ports"
)

// WithdrawalHandler обрабатывает HTTP запросы истории списаний.
type WithdrawalHandler struct {
	balanceService ports.BalanceService
}

// NewWithdrawalHandler создаёт новый обработчик списаний.
func NewWithdrawalHandler(balanceService ports.BalanceService) *WithdrawalHandler {
	return &WithdrawalHandler{
		balanceService: balanceService,
	}
}

// List возвращает историю списаний пользователя.
func (h *WithdrawalHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.FromDomainWithdrawals(withdrawals))
}
