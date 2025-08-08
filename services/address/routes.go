package address

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/services/auth"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
)

type Handler struct {
	store     types.AddressStore
	userStore types.UserStore
}

func NewHandler(store types.AddressStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(auth.WithJWTAuth(h.userStore))

		r.Post("/address", h.handleCreateAddress)
		r.Get("/address", h.handleGetAddress)
		r.Put("/address/{addressID}", h.handleUpdateAddress)
	})
}

func (h *Handler) handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(types.UserKey).(uuid.UUID)

	// parse json
	var input types.CreateAddressPayload

	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate the data
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	address := &types.Address{
		ID:        uuid.New(),
		UserID:    userID,
		Line1:     input.Line1,
		Line2:     input.Line2,
		City:      input.City,
		Country:   input.Country,
		ZipCode:   input.ZipCode,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateAddress(address); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, address)
}

func (h *Handler) handleGetAddress(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(types.UserKey).(uuid.UUID)

	address, err := h.store.GetAddress(userid)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, address)
}

func (h *Handler) handleUpdateAddress(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(types.UserKey).(uuid.UUID)
	var input types.CreateAddressPayload

	// parse json
	if err := utils.ParseJSON(r, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate input
	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get adress id
	addressID := chi.URLParam(r, "addressID")
	parsedID, err := uuid.Parse(addressID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	address := &types.Address{
		ID:      parsedID,
		UserID:  userid,
		Line1:   input.Line1,
		Line2:   input.Line2,
		City:    input.City,
		Country: input.Country,
		ZipCode: input.ZipCode,
	}

	if err := h.store.UpdateAddress(address); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, address)
}
