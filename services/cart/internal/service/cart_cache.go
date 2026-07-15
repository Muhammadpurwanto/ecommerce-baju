package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/model"
	"github.com/redis/go-redis/v9"
)

type CartCacheService interface {
	GetCart(userID string) (*model.Cart, error)
	SetCart(userID string, cart *model.Cart) error
	DeleteCart(userID string) error
}

type cartCacheService struct {
	rdb *redis.Client
}

func NewCartCacheService(rdb *redis.Client) CartCacheService {
	return &cartCacheService{rdb: rdb}
}

func (s *cartCacheService) GetCart(userID string) (*model.Cart, error) {
	ctx := context.Background()
	key := fmt.Sprintf("cart:user:%s", userID)

	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err // redis.Nil jika tidak ada
	}

	var cart model.Cart
	if err := json.Unmarshal([]byte(data), &cart); err != nil {
		return nil, err
	}
	return &cart, nil
}

func (s *cartCacheService) SetCart(userID string, cart *model.Cart) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:user:%s", userID)

	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	// Cache selama 30 menit
	return s.rdb.Set(ctx, key, data, 30*time.Minute).Err()
}

func (s *cartCacheService) DeleteCart(userID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:user:%s", userID)
	return s.rdb.Del(ctx, key).Err()
}
