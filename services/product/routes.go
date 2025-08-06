package product

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/all", h.handleGetProducts)
		r.Get("/{productID}", h.handleGetProduct)
		r.Post("/create", h.handleCreateProduct)
		r.Delete("/delete/{productID}", h.handleDeleteProduct)
		r.Put("/update/{productID}", h.handleUpdateProduct)
	})
}

// @Summary Create a new product
// @Description Add a new product to the catalog
// @Tags Products
// @Accept json
// @Produce json
// @Param product body types.CreateProductPayload true "Product to create"
// @Success 201 {object} types.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/create [post]

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// parse te json from the user input
	var input types.CreateProductPayload

	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// validate the user input
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// create the product and save
	CategoryUUID, err := uuid.Parse(input.CategoryID)
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
		CategoryID:  CategoryUUID,
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

// @Summary Get all products
// @Description Retrieve a list of all products
// @Tags Products
// @Produce json
// @Success 200 {array} types.Product
// @Failure 500 {object} map[string]string
// @Router /products/all [get]

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetAllProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

// @Summary Get product by ID
// @Description Retrieve a single product by its UUID
// @Tags Products
// @Produce json
// @Param productID path string true "Product UUID"
// @Success 200 {object} types.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{productID} [get]

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "productID")
	productUUID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID: %w", err))
		return
	}

	product, err := h.store.GetProductByID(productUUID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

// @Summary Update an existing product
// @Description Modify a product by its UUID
// @Tags Products
// @Accept json
// @Produce json
// @Param productID path string true "Product UUID"
// @Param product body types.CreateProductPayload true "Updated product data"
// @Success 200 {object} types.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/update/{productID} [put]

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var input types.CreateProductPayload

	// Parse JSON body
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate input
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Get product ID from URL
	productIDStr := chi.URLParam(r, "productID")
	productUUID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// Convert CategoryID to UUID
	categoryUUID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid category ID"))
		return
	}

	// Create full product object
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

	// Update in DB
	if err := h.store.UpdateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

// @Summary Delete a product
// @Description Remove a product by its UUID
// @Tags Products
// @Param productID path string true "Product UUID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/delete/{productID} [delete]

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// get the id
	productStr := chi.URLParam(r, "productID")
	productUUID, err := uuid.Parse(productStr)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product id: %w", err))
		return
	}

	if err := h.store.DeleteProduct(productUUID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteNoContent(w)

}
