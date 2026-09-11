package product

import (
	"context"

	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Service interface {
	CreateProduct(ctx context.Context, product models.CreateProductRequest, userID int64) (int64, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service { return &service{repo: repo} }

func (s *service) CreateProduct(ctx context.Context, product models.CreateProductRequest, userID int64) (int64, error) {

	return s.repo.CreateProduct(ctx, product, userID)
}
