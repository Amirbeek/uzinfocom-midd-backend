package product

import (
	"encoding/json"
	"net/http"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/utils"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	id, err := h.svc.CreateProduct(r.Context(), req)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusOK, map[string]int64{"product_id": id})
}
