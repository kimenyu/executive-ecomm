package product

import (
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"net/http"
	"time"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse JSON from request body
	var input types.CreateProductPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the request data
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//  Convert to full Product model
	categoryUUID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product := &types.Product{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Image:       input.Image,
		CategoryID:  categoryUUID,
		Quantity:    input.Quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	//  Save to DB via store
	if err := h.store.CreateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond with the created product
	utils.WriteJSON(w, http.StatusCreated, product)
}
