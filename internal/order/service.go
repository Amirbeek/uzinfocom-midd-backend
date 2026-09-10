package order

import (
	"context"

	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Service interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderID int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID int64) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, order models.Order) error {
	// MASHQ:  POST /orders — bir nechta item'li buyurtma, Idempotency-Key header majburiy (bir xil key bilan qayta yuborilsa, stock ikkinchi marta kamaymasligi kerak)

	return nil
}

func (s *service) GetOrder(ctx context.Context, orderID int64) (*models.Order, error) {
	// MASHQ:  GET /orders/{id} — status: pending → confirmed / cancelled
	return nil, nil
}

func (s *service) CancelOrder(ctx context.Context, orderID int64) error {
	// MASHQ:  POST /orders/{id}/cancel — reserved stock qaytarilishi kerak

	return nil
}
