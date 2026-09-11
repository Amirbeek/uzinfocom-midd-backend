package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/cache"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

type Service interface {
	CreateOrder(ctx context.Context, order models.Order, idempotency string, userID int64) error
	GetOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID int64, userID int64) error
	CancelExpiredOrders(ctx context.Context) error
}

type service struct {
	repo  Repo
	cache *cache.RedisCache
}

func NewService(repo Repo, cache *cache.RedisCache) Service {
	return &service{repo: repo, cache: cache}
}

// Izoh, ochiq muamolar -  hozir redis cache dbdan oldin yozadi, db hato qilsaham biz redisda order bor deb qoladi, keyin agar redisda ishlashda muamo bolsa order ishlamasligi mumkin, Maslahat read-heavy ga yozishni maslahat beraman
func (s *service) CreateOrder(ctx context.Context, order models.Order, idempotency string, userId int64) error {
	// MASHQ:  POST /orders — bir nechta item'li buyurtma, Idempotency-Key header majburiy (bir xil key bilan qayta yuborilsa, stock ikkinchi marta kamaymasligi kerak)
	key := fmt.Sprintf("order:%d:%s", userId, idempotency)
	// order:{userID}:{idempotency} qilib olamiz bu unique bolishi kerak. keyin uni createOrder qiilib redisdan qidiramiz, agar topilsa cashed qilib return qilamiz bu ram da saqlanganligi tufayli tezroq va diska qaraganda ancha tez hisoblanadi
	if cached, err := s.cache.GetOrder(ctx, key); err == nil && cached != nil {
		return nil
	}
	// agar redisda topilmasa db dan olib kelamiz va redis ga set qilamiz
	if err := s.repo.CreateOrder(ctx, order, idempotency, userId); err != nil {
		return err
	}
	_ = s.cache.SetOrder(ctx, key, &order, time.Minute)

	return nil
}

func (s *service) GetOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error) {
	// MASHQ:  GET /orders/{id} — status: pending → confirmed / cancelled
	key := fmt.Sprintf("order:%d:%d", userID, orderID)
	// order:{userID}:{orderID} qilib olamiz bu unique bolishi kerak. keyin uni getOrder qiilib redisdan qidiramiz, agar topilsa cashed qilib return qilamiz bu ram da saqlanganligi tufayli tezroq va diska qaraganda ancha tez hisoblanadi
	if cached, err := s.cache.GetOrder(ctx, key); err == nil && cached != nil {
		return cached, nil
	}
	// agar redisda topilmasa db dan olib kelamiz va redis ga set qilamiz
	order, err := s.repo.GetOrder(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetOrder(ctx, key, order, time.Minute)

	return order, nil
}

func (s *service) CancelOrder(ctx context.Context, orderID int64, userID int64) error {
	// MASHQ:  POST /orders/{id}/cancel — reserved stock qaytarilishi kerak
	// bu yerda radis cache dan orderni delete qilamiz, chunki ode
	if err := s.repo.CancelOrder(ctx, orderID, userID); err != nil {
		return err
	}
	_ = s.cache.DeleteOrder(ctx, fmt.Sprintf("order:%d:%d", userID, orderID))

	return nil
}

func (s *service) CancelExpiredOrders(ctx context.Context) error {
	// MASHQ:  Crone job har 15 minutda backroundda 15 minutdan otkan pending orderlarni cancel qilib stock quantityni qaytarib qoyadi
	return s.repo.CancelExpiredOrders(ctx)
}
