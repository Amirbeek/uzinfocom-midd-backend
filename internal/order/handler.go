package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return err
	}
	return h.svc.CreateOrder(r.Context(), order)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) error {
	orderID, err := orderIDFromURL(r)
	if err != nil {
		return err
	}

	order, err := h.svc.GetOrder(r.Context(), orderID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(order)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) error {
	orderID, err := orderIDFromURL(r)
	if err != nil {
		return err
	}

	return h.svc.CancelOrder(r.Context(), orderID)
}

func orderIDFromURL(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}
