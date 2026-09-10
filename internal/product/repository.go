package product

import (
	"context"
	"database/sql"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Repo interface {
	CreateProduct(ctx context.Context, product models.Product, userId int64) (int64, error)
}

type repo struct{ db *sql.DB }

func NewRepo(db database.Service) Repo {
	return &repo{db: db.DB()}
}

func (r *repo) CreateProduct(ctx context.Context, product models.Product, userId int64) (int64, error) {
	// MASHQ:  POST /products — mahsulot yaratish (name, price, stock_quantity)

	q := `
		INSERT INTO products(
		user_id,
		name, 
		price,
		stock_quantity
		)VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx, q, userId, product.Name, product.Price, product.StockQuantity,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}
