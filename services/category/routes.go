package category

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"net/http"
	"time"
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
		r.Get("/{productID}", h.handleGetCategory)
		r.Post("/", h.handleCreateCategory)

	})
}

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

// get categories handler
func (h *Handler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.GetCategories()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, categories)
}

// get singe category handler
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
