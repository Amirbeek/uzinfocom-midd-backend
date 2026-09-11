package user

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

// Register godoc
//
//	@Summary	Login a user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body	models.LoginRequest	true	"Login request"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	utils.ErrorResponse
//	@Failure	401	{object}	utils.ErrorResponse
//	@Router		/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	token, err := h.svc.Login(r.Context(), req)
	if err != nil {
		utils.UnauthorizedError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

// Register godoc
//
//	@Summary	Register a new user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body	models.RegisterRequest	true	"Register request"
//	@Success	200	{object}	map[string]string
//	@Failure	400	{object}	utils.ErrorResponse
//	@Failure	401	{object}	utils.ErrorResponse
//	@Router		/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	token, err := h.svc.Register(r.Context(), req)
	if err != nil {
		utils.UnauthorizedError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
