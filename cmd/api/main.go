package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/Amirbeek/uzinfocom-midd-backend/docs"
)

const version = "0.1.0"

//	@title			Uzinfocom Backend API
//	@version		0.1.0
//	@description	Backend API for the Uzinfocom middle backend developer task.
//	@host			localhost:8080
//	@BasePath		/v1

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", healthHandler)
	mux.Handle("GET /v1/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/v1/swagger/doc.json"),
	))

	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

type healthResponse struct {
	Status  string `json:"status"  example:"ok"`
	Version string `json:"version" example:"0.1.0"`
}

// healthHandler godoc
//
//	@Summary	Health check
//	@Tags		ops
//	@Produce	json
//	@Success	200	{object}	main.healthResponse
//	@Router		/health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(healthResponse{
		Status:  "ok",
		Version: version,
	})
}
