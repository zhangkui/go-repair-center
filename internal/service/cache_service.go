package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	redisv9 "github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redisv9.Client
	locker *redislock.Client
	ttl    time.Duration
}

func NewCacheService(client *redisv9.Client, ttl time.Duration) *CacheService {
	return &CacheService{
		client: client,
		locker: redislock.New(client),
		ttl:    ttl,
	}
}

func (s *CacheService) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = s.ttl
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, payload, ttl).Err()
}

func (s *CacheService) GetJSON(ctx context.Context, key string, target any) (bool, error) {
	payload, err := s.client.Get(ctx, key).Bytes()
	if err == redisv9.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return false, err
	}
	return true, nil
}

func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *CacheService) DeleteMany(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}

func (s *CacheService) RememberJSON(ctx context.Context, key string, ttl time.Duration, target any, loader func() (any, error)) error {
	if ok, err := s.GetJSON(ctx, key, target); err != nil {
		return err
	} else if ok {
		return nil
	}

	value, err := loader()
	if err != nil {
		return err
	}
	if err := s.SetJSON(ctx, key, value, ttl); err != nil {
		return err
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, target)
}

func (s *CacheService) Lock(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	if ttl <= 0 {
		ttl = 15 * time.Second
	}
	lock, err := s.locker.Obtain(ctx, key, ttl, nil)
	if err == redislock.ErrNotObtained {
		return nil, fmt.Errorf("lock busy: %s", key)
	}
	if err != nil {
		return nil, err
	}
	return lock, nil
}

func (s *CacheService) WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	lock, err := s.Lock(ctx, key, ttl)
	if err != nil {
		return err
	}
	operationErr := fn()
	releaseErr := lock.Release(ctx)
	return CombineLockErrors(operationErr, releaseErr)
}

func (s *CacheService) BuildKey(parts ...string) string {
	result := "repair-center"
	for _, part := range parts {
		if part == "" {
			continue
		}
		result += ":" + part
	}
	return result
}

func (s *CacheService) TTL() time.Duration {
	return s.ttl
}
