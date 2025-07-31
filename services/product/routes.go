package product

import (
	"fmt"
	"github.com/go-chi/chi/v5"
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

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.handleGetProducts)
		r.Get("/{productID}", h.handleGetProduct)
		r.Post("/", h.handleCreateProduct)
	})
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var input types.CreateProductPayload

	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

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

	if err := h.store.CreateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetAllProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "productID")
	productUUID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.store.GetProductByID(productUUID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	//  Parse productID from URL
	productIDStr := chi.URLParam(r, "productID")
	productUUID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID: %w", err))
		return
	}

	//  Parse JSON body
	var input types.CreateProductPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		return
	}

	//  Validate the input
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %w", err))
		return
	}

	// Parse category UUID
	categoryUUID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid category ID: %w", err))
		return
	}

	// Create updated product struct
	product := &types.Product{
		ID:          productUUID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Image:       input.Image,
		CategoryID:  categoryUUID,
		Quantity:    input.Quantity,
		UpdatedAt:   time.Now(),
	}

	//  Update in DB
	if err := h.store.UpdateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update product: %w", err))
		return
	}

	//  Respond with updated product
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	productIDStr := chi.URLParam(r, "productID")
	productUUID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID: %w", err))
		return
	}

	// Delete product using the store
	if err := h.store.DeleteProduct(productUUID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete product: %w", err))
		return
	}

	// Return 204 No Content
	utils.WriteNoContent(w)
}
