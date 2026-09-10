package main

import (
	_ "github.com/Amirbeek/uzinfocom-midd-backend/docs"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/env"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/server"
)

const version = "0.1.0"

// @title			Uzinfocom Backend API
// @version		0.1.0
// @description	Backend API for the Uzinfocom middle backend developer task.
// @host			localhost:8080
// @BasePath		/v1
func main() {
	cfg := server.Config{
		Addr: env.GetString("APP_ADDR", ":8080"),
		JWT: server.JWTConfig{
			Secret: env.GetString("JWT_SECRET", "change-me-in-production"),
			Aud:    env.GetString("JWT_AUD", "uzinfocom"),
			Iss:    env.GetString("JWT_ISS", "uzinfocom"),
		},
	}

	app := server.NewApplication(cfg)

	mux := app.ServeHTTP()

	apiServer := app.Run(mux)

	done := make(chan bool, 1)

	go gracefulShutdown(apiServer, done)

	<-done
}
