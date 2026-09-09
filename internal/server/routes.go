package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) ServeHTTP() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // panik bolgan holatlarda serverni qayta uygotadi va 500 hato tashaydi
	r.Use(middleware.RequestID) // har bitta request uchun unique id beradi
	r.Use(middleware.RealIP)    // haqiqiy ip manzilini olish uchun ishlatiladi

	return r
}
