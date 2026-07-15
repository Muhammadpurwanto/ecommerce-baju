package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenCacheService interface {
	BlacklistToken(token string, expiration time.Duration) error
	IsTokenBlacklisted(token string) bool
}

type tokenCacheService struct {
	rdb *redis.Client
}

func NewTokenCacheService(rdb *redis.Client) TokenCacheService {
	return &tokenCacheService{rdb: rdb}
}

func (s *tokenCacheService) BlacklistToken(token string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", token)
	return s.rdb.Set(ctx, key, "revoked", expiration).Err()
}

func (s *tokenCacheService) IsTokenBlacklisted(token string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", token)
	val, err := s.rdb.Get(ctx, key).Result()
	return err == nil && val == "revoked"
}
