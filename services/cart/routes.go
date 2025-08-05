package cart

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/services/auth"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"net/http"
	"time"
)

type Handler struct {
	store     types.CartStore
	userStore types.UserStore
}

func NewHandler(store types.CartStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(h.userStore))

		r.Post("/products/{productID}/cart", h.handleAddItemToCart)
	})
}

func (h *Handler) handleAddItemToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	// Get product ID from URL
	productStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var input types.AddToCartPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check if cart exists for the user
	cart, err := h.store.GetCartByUserID(userID)
	if err == sql.ErrNoRows {
		// Create new cart
		newCart := &types.Cart{
			ID:        uuid.New(),
			UserID:    userID,
			CreatedAt: time.Now(),
		}
		if err := h.store.CreateCart(newCart); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		cart = newCart
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Add item to cart
	cartItem := &types.CartItem{
		ID:        uuid.New(),
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  input.Quantity,
		CreatedAt: time.Now(),
	}
	if err := h.store.AddCartItem(cartItem); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, cartItem)
}

func (h *Handler) handleGetCartItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	cart, err := h.store.GetCartByUserID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	items, err := h.store.GetCartItems(cart.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, items)
}
