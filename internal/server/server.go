package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/auth"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/order"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/product"
	services "github.com/Amirbeek/uzinfocom-midd-backend/internal/service"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/user"
)

type Application struct {
	Config         Config
	db             database.Service
	Services       *services.Services
	Auth           *auth.JWTAuthenticator
	productHandler *product.Handler
	orderHandler   *order.Handler
	user           *user.Handler
}

func NewApplication(config Config) *Application {
	db := database.New()

	// Business Logic Layer (db -> Repo -> Service)
	svcs := services.NewServices(db)

	jwtAuth := auth.NewJWTAuthenticator(config.JWT.Secret, config.JWT.Aud, config.JWT.Iss)

	return &Application{
		Config:         config,
		db:             db,
		Services:       svcs,
		Auth:           jwtAuth,
		productHandler: product.NewHandler(svcs.Product),
		orderHandler:   order.NewHandler(svcs.Order),
	}
}

type Config struct {
	Addr string
	JWT  JWTConfig
}

type JWTConfig struct {
	Secret string
	Aud    string
	Iss    string
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
