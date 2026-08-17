// Package cache Redis 缓存封装：验证码、限流计数
package cache

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis 缓存客户端
type Redis struct {
	client *redis.Client
}

// NewRedis 创建 Redis 客户端并 Ping
func NewRedis(addr, password string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return &Redis{client: client}, nil
}

// Set 写入带过期时间的键
func (r *Redis) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get 读取键值
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// GenerateCode 生成 6 位数字验证码并写入 Redis
func (r *Redis) GenerateCode(ctx context.Context, key string, ttl time.Duration) (string, error) {
	code, err := randomCode()
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}
	if err := r.Set(ctx, key, code, ttl); err != nil {
		return "", fmt.Errorf("set code: %w", err)
	}
	return code, nil
}

func randomCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
