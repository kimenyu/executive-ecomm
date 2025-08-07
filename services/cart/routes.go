package cart

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/services/auth"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
)

type Handler struct {
	store     types.CartStore
	userStore types.UserStore
}

func NewHandler(store types.CartStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(h.userStore))

		r.Post("/products/{productID}/cart", h.handleAddItemToCart)
		r.Get("/cart/my/items", h.handleGetCartItems)
	})
}

// @Summary Add product to cart
// @Description Add a product to the authenticated user's cart
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param productID path string true "Product UUID"
// @Param payload body types.AddToCartPayload true "Cart item payload"
// @Success 201 {object} types.CartItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{productID}/cart [post]

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

// @Summary Get my cart items
// @Description Retrieve items in the authenticated user's cart
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {array} types.CartItem
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/my/items [get]

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
