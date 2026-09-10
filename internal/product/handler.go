package product

import (
	"encoding/json"
	"errors"
	"net/http"

	mw "github.com/Amirbeek/uzinfocom-midd-backend/internal/middleware"
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
	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok {
		utils.UnauthorizedError(
			w,
			r,
			errors.New("user id not found"),
		)
		return
	}

	var req models.Product

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	id, err := h.svc.CreateProduct(r.Context(), req, userID)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusCreated, map[string]int64{"product_id": id})
}
