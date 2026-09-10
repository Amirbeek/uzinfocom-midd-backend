package service

import (
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/order"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/product"
)

// Services — business logic qatlami. Har bir domen uchun
// db -> Repo -> Service zanjiri shu yerda quriladi.
type Services struct {
	Product product.Service
	Order   order.Service
}

func NewServices(db database.Service) *Services {
	return &Services{
		Product: product.NewService(product.NewRepo(db)),
		Order:   order.NewService(order.NewRepo(db)),
	}
}
