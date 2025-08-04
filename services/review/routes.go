package review

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/services/auth"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"net/http"
	"time"
)

type Handler struct {
	store     types.ReviewStore
	userStore types.UserStore
}

func NewHandler(store types.ReviewStore) *Handler {
	return &Handler{store: store}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(h.userStore)) // injects user ID into context

		r.Post("/products/{productID}/reviews", h.handleCreateReview)
		r.Get("/products/{productID}/reviews", h.handleGetReviewsByProduct)
		r.Get("/reviews/{id}", h.handleGetReviewByID)
		r.Put("/reviews/{id}", h.handleUpdateReview)
		r.Delete("/reviews/{id}", h.handleDeleteReview)
	})
}

func (h *Handler) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	// Get product ID from URL
	productStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var input types.CreateReviewPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	review := &types.Review{
		ID:        uuid.New(),
		ProductID: productID,
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateReview(review); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, review)
}

func (h *Handler) handleDeleteReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	reviewIDStr := chi.URLParam(r, "id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.DeleteReview(reviewID, userID); err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}

	utils.WriteNoContent(w)
}

func (h *Handler) handleGetReviewByID(w http.ResponseWriter, r *http.Request) {
	reviewID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	review, err := h.store.GetReviewByID(reviewID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, review)
}

func (h *Handler) handleGetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	reviews, err := h.store.GetReviewsByProduct(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, reviews)
}

func (h *Handler) handleUpdateReview(w http.ResponseWriter, r *http.Request) {
	// Parse review ID from URL
	reviewIDStr := chi.URLParam(r, "id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid review ID"))
		return
	}

	// Extract user ID from JWT-authenticated context
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	// Fetch original review to ensure user is the owner and get product ID
	existingReview, err := h.store.GetReviewByID(reviewID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("review not found"))
		return
	}

	if existingReview.UserID != userID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized to update this review"))
		return
	}

	// Parse incoming review payload
	var input types.CreateReviewPayload
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Construct updated review
	updatedReview := &types.Review{
		ID:        reviewID,
		ProductID: existingReview.ProductID,
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		CreatedAt: existingReview.CreatedAt,
	}

	// Save changes
	if err := h.store.UpdateReview(updatedReview); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedReview)
}
