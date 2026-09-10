package product

import (
	"context"
	"database/sql"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Repo interface {
	CreateProduct(ctx context.Context, product models.Product) (int64, error)
}

type repo struct{ db *sql.DB }

func NewRepo(db database.Service) Repo {
	return &repo{db: db.DB()}
}

func (r *repo) CreateProduct(ctx context.Context, product models.Product) (int64, error) {
	// MASHQ:  POST /products — mahsulot yaratish (name, price, stock_quantity)

	return 0, nil
}
