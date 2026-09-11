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

// CreateProduct godoc
//
//	@Summary	Mahsulot yaratish
//	@Tags		products
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		request	body	models.CreateProductRequest	true	"Product"
//	@Success	201	{object}	map[string]int64
//	@Failure	400	{object}	utils.ErrorResponse
//	@Failure	401	{object}	utils.ErrorResponse
//	@Router		/products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok {
		utils.UnauthorizedError(w, r, errors.New("user id not found"))
		return
	}

	var req models.CreateProductRequest

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
