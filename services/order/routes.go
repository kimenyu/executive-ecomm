package order

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/services/auth"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
)

type Handler struct {
	store        types.OrderStore
	userStore    types.UserStore
	addressStore types.AddressStore
}

func NewHandler(store types.OrderStore, userStore types.UserStore, addressStore types.AddressStore) *Handler {
	return &Handler{store: store, userStore: userStore, addressStore: addressStore}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(h.userStore))
		r.Post("/orders", h.handleCreateOrder)
		r.Get("/orders", h.handleGetOrdersByUser)
	})
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	var input types.CreateOrderPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate items
	if len(input.Items) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("order must contain at least one item"))
		return
	}

	// Calculate total
	var total float64
	for _, item := range input.Items {
		if item.Quantity <= 0 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("quantity for product %s must be greater than 0", item.ProductID))
			return
		}
		if item.Price <= 0 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("price for product %s must be greater than 0", item.ProductID))
			return
		}
		total += float64(item.Quantity) * item.Price
	}

	address, err := h.addressStore.GetAddress(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// Create order
	order := &types.Order{
		ID:        uuid.New(),
		UserID:    userID,
		AddressID: address.ID,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateOrder(order); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Add order items
	for _, item := range input.Items {
		orderItem := &types.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if err := h.store.AddOrderItem(orderItem); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.WriteJSON(w, http.StatusCreated, order)
}

func (h *Handler) handleGetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	orders, err := h.store.GetOrdersByUser(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}
