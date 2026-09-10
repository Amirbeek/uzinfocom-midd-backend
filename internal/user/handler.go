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
