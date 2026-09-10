package server

import (
	"net/http"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/utils"
)

type healthResponse struct {
	Status string `json:"status" example:"ok"`
}

// healthHandler godoc
//
//	@Summary	Health check
//	@Tags		ops
//	@Produce	json
//	@Success	200	{object}	server.healthResponse
//	@Router		/health [get]
func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	_ = utils.WriteJson(w, http.StatusOK, healthResponse{Status: "ok"})
}
