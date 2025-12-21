package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/arvaliullin/gophermart/internal/api/http/dto"
	"github.com/arvaliullin/gophermart/internal/api/http/middleware"
	"github.com/arvaliullin/gophermart/internal/core/domain"
	"github.com/arvaliullin/gophermart/internal/core/ports"
)

// OrderHandler обрабатывает HTTP запросы управления заказами.
type OrderHandler struct {
	orderService ports.OrderService
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Submit обрабатывает загрузку номера заказа.
func (h *OrderHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		http.Error(w, ErrInvalidRequestFormat.Error(), http.StatusBadRequest)
		return
	}

	alreadyExists, err := h.orderService.SubmitOrder(r.Context(), userID, orderNumber)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidOrderNumber) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, domain.ErrOrderBelongsToOther) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if alreadyExists {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// List возвращает список заказов пользователя.
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetUserOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.FromDomainOrders(orders))
}
