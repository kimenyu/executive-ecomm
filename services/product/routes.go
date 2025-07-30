package product

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// parse the json
	var p types.CreateProductPayload

	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate struct data
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Convert to full Product model
	categoryUUID, err := uuid.Parse(p.CategoryID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product := &types.Product{
		ID:          uuid.New(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
		CategoryID:  categoryUUID,
		Quantity:    p.Quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// store te product
	if err := h.store.CreateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJSON(w, http.StatusCreated, product)
}

func (h *Handler) handleGetProductsById(w http.ResponseWriter, r *http.Request) {
	// extract the id from the urls path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("invalid URL forma, expected /products/{id}"))
		return
	}

	// parse uuid
	id, err := uuid.Parse(parts[2])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
	}

	// fetch from the db using store
	product, err := h.store.GetProductByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	// return product as JSON
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetAllProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

// update product handler
