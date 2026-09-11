package main

import (
	"time"

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
			TTL:    24 * time.Hour,
		},
	}

	app := server.NewApplication(cfg)

	// Orqada expired bo'lgan pending buyurtmalarni bekor qiladi.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.StartExpiredOrderCancellation(ctx, time.Duration(env.GetInt("ORDER_CLEANUP_INTERVAL_SECONDS", 60))*time.Second)

	mux := app.ServeHTTP()

	apiServer := app.Run(mux)

	done := make(chan bool, 1)

	go gracefulShutdown(apiServer, done)

	<-done
}
