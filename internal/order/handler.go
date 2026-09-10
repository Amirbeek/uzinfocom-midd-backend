package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	mw "github.com/Amirbeek/uzinfocom-midd-backend/internal/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/utils"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")

	if key == "" {
		utils.BadRequestError(w, r, errors.New("The operation not supported"))
		return
	}
	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("user id not found"))
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.svc.CreateOrder(r.Context(), order, key, userID); err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusCreated, order)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := orderIDFromURL(r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("user id not found"))
		return
	}

	order, err := h.svc.GetOrder(r.Context(), orderID, userID)
	if err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusOK, order)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := orderIDFromURL(r)

	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("user id not found"))
		return
	}

	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.svc.CancelOrder(r.Context(), orderID, userID); err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func orderIDFromURL(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}
