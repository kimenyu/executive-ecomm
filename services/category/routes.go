package category

import (
	"github.com/kimenyu/executive/types"
	"net/http"
)

type Handler struct {
	store types.CategoryStore
}

func NewHandler(store types.CategoryStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleCreateCategories(w http.ResponseWriter, r *http.Request) {
	
}
