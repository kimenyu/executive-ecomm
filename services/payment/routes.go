package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
)

type Handler struct {
	store      *Store
	orderStore types.OrderStore
}

func NewHandler(store *Store, orderStore types.OrderStore) *Handler {
	return &Handler{store: store, orderStore: orderStore}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/payments/confirm", h.handleConfirm)
}

// payload matches what Node sends
type confirmPayload struct {
	OrderID         string          `json:"order_id"`
	Status          string          `json:"status"`
	Amount          float64         `json:"amount"`
	Provider        string          `json:"provider"`
	CheckoutRequest string          `json:"checkout_request_id"`
	MerchantRequest string          `json:"merchant_request_id"`
	MpesaReceipt    string          `json:"mpesa_receipt"`
	Phone           string          `json:"phone"`
	Raw             json.RawMessage `json:"raw"`
}

func (h *Handler) handleConfirm(w http.ResponseWriter, r *http.Request) {
	// verify secret
	expected := os.Getenv("NODE_NOTIFY_SECRET")
	if r.Header.Get("X-Node-Notify-Secret") != expected {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	var p confirmPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	orderUUID, err := uuid.Parse(p.OrderID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order_id"))
		return
	}

	order, err := h.orderStore.GetOrderWithItemsByID(orderUUID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("order not found"))
		return
	}

	if p.Amount != order.Order.Total {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payment amount %.2f does not match order total %.2f", p.Amount, order.Order.Total))
		return
	}

	// create payment record
	pay := &types.Payment{
		ID:                uuid.New(),
		OrderID:           order.Order.ID,
		Amount:            p.Amount,
		Provider:          p.Provider,
		Status:            p.Status,
		CheckoutRequestID: p.CheckoutRequest,
		MerchantRequestID: p.MerchantRequest,
		MpesaReceipt:      p.MpesaReceipt,
		Phone:             p.Phone,
		Metadata:          p.Raw,
		CreatedAt:         time.Now(),
	}

	if err := h.store.CreatePayment(pay); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// update order state if success
	if p.Status == "success" {
		if err := h.orderStore.UpdateOrderStatus(order.Order.ID, "paid"); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	fmt.Printf("Received confirmPayload.OrderID: %v\n", p.OrderID)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
