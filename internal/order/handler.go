package order

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.svc.CreateOrder(r.Context(), order); err != nil {
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

	order, err := h.svc.GetOrder(r.Context(), orderID)
	if err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusOK, order)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := orderIDFromURL(r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.svc.CancelOrder(r.Context(), orderID); err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func orderIDFromURL(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}
