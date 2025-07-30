package category

import (
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

func (h *Handler) handleCreateCategories(w http.ResponseWriter, r *http.Request) {

	// parse the category
	var c types.CreateCategoryPayload

	if err := utils.ParseJSON(r, &c); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the incoming json
	if err := utils.Validate.Struct(c); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//create the category and save
	category := &types.Category{
		ID:        uuid.New(),
		Name:      c.Name,
		CreatedAt: time.Now(),
	}

	// saving
	if err := h.store.CreateCategory(category); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJSON(w, http.StatusCreated, category)
}

// get all categories
func (h *Handler) handleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.GetCategories()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJSON(w, http.StatusOK, categories)
}
