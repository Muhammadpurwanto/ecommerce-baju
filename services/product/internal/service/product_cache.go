package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
	"github.com/redis/go-redis/v9"
)

type ProductCacheService interface {
	GetProducts(page, perPage int) ([]model.Product, error)
	SetProducts(page, perPage int, products []model.Product) error
	GetProduct(slug string) (*model.Product, error)
	SetProduct(slug string, product *model.Product) error
	InvalidateAll() error
}

type productCacheService struct {
	rdb *redis.Client
}

func NewProductCacheService(rdb *redis.Client) ProductCacheService {
	return &productCacheService{rdb: rdb}
}

func (s *productCacheService) GetProducts(page, perPage int) ([]model.Product, error) {
	ctx := context.Background()
	key := fmt.Sprintf("products:list:%d:%d", page, perPage)

	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var products []model.Product
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productCacheService) SetProducts(page, perPage int, products []model.Product) error {
	ctx := context.Background()
	key := fmt.Sprintf("products:list:%d:%d", page, perPage)

	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, key, data, 10*time.Minute).Err()
}

func (s *productCacheService) GetProduct(slug string) (*model.Product, error) {
	ctx := context.Background()
	key := fmt.Sprintf("products:slug:%s", slug)

	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var product model.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *productCacheService) SetProduct(slug string, product *model.Product) error {
	ctx := context.Background()
	key := fmt.Sprintf("products:slug:%s", slug)

	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, key, data, 10*time.Minute).Err()
}

func (s *productCacheService) InvalidateAll() error {
	ctx := context.Background()
	// Hapus semua cache product
	iter := s.rdb.Scan(ctx, 0, "products:*", 100).Iterator()
	for iter.Next(ctx) {
		s.rdb.Del(ctx, iter.Val())
	}
	return iter.Err()
}
