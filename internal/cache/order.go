package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

func (r *RedisCache) GetOrder(ctx context.Context, key string) (*models.Order, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *RedisCache) SetOrder(ctx context.Context, key string, order *models.Order, expiration time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) DeleteOrder(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
