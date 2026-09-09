package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	services "github.com/Amirbeek/uzinfocom-midd-backend/internal/service"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/store"
)

type Application struct {
	Config   Config
	db       database.Service
	Services *services.Services
}

func NewApplication(config Config) *Application {
	db := database.New()

	// database layer
	st := store.NewStore(db)

	// Business Logic Layer
	svcs := services.NewServices(st)

	return &Application{
		Config:   config,
		db:       db,
		Services: svcs,
	}
}

type Config struct {
	Addr string
}

func (app *Application) DB() database.Service {
	return app.db
}

func (app *Application) Run(mux http.Handler) *http.Server {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("Listening on %s (HTTP)", app.Config.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	return server
}
