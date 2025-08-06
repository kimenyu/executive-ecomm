package category

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
)

type Handler struct {
	store types.CategoryStore
}

func NewHandler(store types.CategoryStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", h.handleGetCategories)
		r.Get("/{id}", h.handleGetCategory)
		r.Post("/", h.handleCreateCategory)

	})
}

// @Summary Create a new category
// @Description Add a new category to group products
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body types.CreateCategoryPayload true "Category to create"
// @Success 201 {object} types.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/ [post]

func (h *Handler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	//parse the json
	var input types.CreateCategoryPayload

	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	category := &types.Category{
		ID:        uuid.New(),
		Name:      input.Name,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateCategory(category); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, category)
}

// @Summary Get all categories
// @Description Retrieve all available product categories
// @Tags Categories
// @Produce json
// @Success 200 {array} types.Category
// @Failure 500 {object} map[string]string
// @Router /categories/ [get]

func (h *Handler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.GetCategories()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, categories)
}

// @Summary Get category by ID
// @Description Retrieve a single category by its UUID
// @Tags Categories
// @Produce json
// @Param id path string true "Category UUID"
// @Success 200 {object} types.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [get]

func (h *Handler) handleGetCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryUUID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	category, err := h.store.GetCategoryById(categoryUUID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, category)
}
