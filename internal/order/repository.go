package order

import (
	"context"
	"database/sql"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Repo interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderID int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID int64) error
}

type repo struct{ db *sql.DB }

func NewRepo(db database.Service) Repo {
	return &repo{db: db.DB()}
}

func (r *repo) CreateOrder(ctx context.Context, order models.Order) error {
	//Mashq: POST /orders — bir nechta item'li buyurtma, Idempotency-Key header majburiy (bir xil key bilan qayta yuborilsa, stock ikkinchi marta kamaymasligi kerak)

	return nil
}

func (r *repo) GetOrder(ctx context.Context, orderID int64) (*models.Order, error) {
	// MASHQ:  GET /orders/{id} — status: pending → confirmed / cancelled

	return nil, nil
}

func (r *repo) CancelOrder(ctx context.Context, orderID int64) error {
	// POST /orders/{id}/cancel — reserved stock qaytarilishi kerak

	return nil
}
