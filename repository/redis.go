package repository

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	rdb *redis.Client
	ctx *context.Context
}

func NewRedisRepository() *RedisRepository {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})
	return &RedisRepository{
		rdb: rdb,
		ctx: &ctx,
	}
}

func (re *RedisRepository) Set(key, value string) error {
	if err := re.rdb.Set(*re.ctx, key, value, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (re *RedisRepository) Get(key string) (string, error) {
	result, err := re.rdb.Get(*re.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
