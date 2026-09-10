package store

import (
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/order"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/product"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/user"
)

// Store — ma'lumot qatlami. Har bir domenning Repo implementatsiyasi
type Store struct {
	Product product.Repo
	Order   order.Repo
	User    user.Repo
}

func NewStore(db database.Service) *Store {
	return &Store{
		Product: product.NewRepo(db),
		Order:   order.NewRepo(db),
		User:    user.NewRepo(db),
	}
}
