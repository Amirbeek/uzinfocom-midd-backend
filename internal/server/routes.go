package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	mw "github.com/Amirbeek/uzinfocom-midd-backend/internal/middleware"
)

func (app *Application) ServeHTTP() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // panik bolgan holatlarda serverni qayta uygotadi va 500 hato tashaydi
	r.Use(middleware.RequestID) // har bitta request uchun unique id beradi
	r.Use(middleware.RealIP)    // haqiqiy ip manzilini olish uchun ishlatiladi

	r.Post("/login", app.user.Login)
	r.Post("/register", app.user.Register)

	r.Route("/v1", func(r chi.Router) {
		// ochiq route lar
		r.Get("/health", app.healthHandler)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/v1/swagger/doc.json"),
		))

		// jwt himoyalangan guruh
		r.Group(func(r chi.Router) {
			r.Use(mw.RequireAuth(app.Auth))

			r.Route("/products", func(r chi.Router) {
				r.Post("/", app.productHandler.CreateProduct)
			})

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", app.orderHandler.CreateOrder)
				r.Get("/{id}", app.orderHandler.GetOrder)
				r.Post("/{id}/cancel", app.orderHandler.CancelOrder)
			})
		})
	})

	return r
}
