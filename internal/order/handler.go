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

// CreateOrder godoc
//
//		@Summary	Order yaratish
//		@Tags		orders
//		@Security	BearerAuth
//	   @Param  Idempotency-Key header  string  true    "idempotentlik kaliti"
//		@Accept		json
//		@Produce=json
//		@Param		request	body	models.CreateOrderRequest	true	"Buyurtma itemlari"
//		@Success	201	{object}	map[string]int64
//		@Failure	400	{object}	utils.ErrorResponse
//		@Failure	401	{object}	utils.ErrorResponse
//		@Router		/orders [post]
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

	var order models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	orderID, err := h.svc.CreateOrder(r.Context(), order, key, userID)
	if err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusCreated, map[string]int64{"order_id": orderID})
}

// GetOrder godoc
//
//	@Summary	Orderni olish
//	@Tags		orders
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path	int	true	"Order id"
//	@Success	200	{object}	models.Order
//	@Failure	404	{object}	utils.ErrorResponse
//	@Failure	400	{object}	utils.ErrorResponse
//	@Failure	401	{object}	utils.ErrorResponse
//	@Router		/orders/{id} [get]
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
	// Bug tuzatildi, oldin sql no rows in result set, degan hato chiqardi hozir bu qism bilan  403 qaytaradi
	if errors.Is(err, ErrorNotFound) {
		utils.NotFoundError(w, r, err)
		return
	}
	if err != nil {
		utils.IntervalServerError(w, r, err)
		return
	}

	_ = utils.WriteJson(w, http.StatusOK, order)
}

// CancelOrder godoc
//
//	@Summary	Orderni bekor qilish
//	@Tags		orders
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path	int	true	"Order id"
//	@Success	204	"Bekor qilindi, reserved stock qaytarildi"
//	@Failure	400	{object}	utils.ErrorResponse
//	@Failure	401	{object}	utils.ErrorResponse
//	@Router		/orders/{id}/cancel [post]
func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.CancelOrder(r.Context(), orderID, userID); err != nil {
		if errors.Is(err, ErrorNotFound) {
			utils.NotFoundError(w, r, err)
			return
		}
		utils.IntervalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func orderIDFromURL(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}
