package service

import (
	"time"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/auth"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/cache"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/order"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/product"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/store"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/user"
)

// Services — biznes mantiq qatlami. Repo'larni Store'dan tayyor holda oladi.
type Services struct {
	Product product.Service
	Order   order.Service
	User    user.Service
}

func NewServices(s *store.Store, a auth.Authenticator, ttl time.Duration, c *cache.RedisCache) *Services {
	return &Services{
		Product: product.NewService(s.Product),
		Order:   order.NewService(s.Order, c),
		User:    user.NewService(s.User, a, ttl),
	}
}
