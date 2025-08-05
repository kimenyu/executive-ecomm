package order

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/services/auth"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"net/http"
	"time"
)

type Handler struct {
	store     types.OrderStore
	userStore types.UserStore
}

func NewHandler(store types.OrderStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(h.userStore))

		r.Post("/orders", h.handleCreateOrder)
		r.Get("/orders", h.handleGetOrders)
	})
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	var input types.CreateOrderPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	order := &types.Order{
		ID:        uuid.New(),
		UserID:    userID,
		Total:     input.Total,
		Status:    "pending",
		AddressID: input.AddressID,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateOrder(order); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

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

func (h *Handler) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	orders, err := h.store.GetOrdersByUser(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}
