package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	rconf := struct {
		IP       string
		Port     string
		Password string
	}{
		IP:       os.Getenv("REDIS_IP"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	fmt.Println(rconf)
	client := redis.NewClient(&redis.Options{
		Addr:     rconf.IP + ":" + rconf.Port,
		Password: rconf.Password,
		DB:       0,
	})
	// 接続確認
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Expire(key string, expiration time.Duration) error
}

type redisRepository struct {
	cli *redis.Client
	ctx context.Context
}

func NewRedisClient(cli *redis.Client, ctx context.Context) RedisRepository {
	return &redisRepository{
		cli: cli,
		ctx: ctx,
	}
}

func (r *redisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.cli.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "redis can't set: %w")
	}
	return nil
}

func (r *redisRepository) Get(key string) (string, error) {
	s, err := r.cli.Get(r.ctx, key).Result()
	if err != nil {
		return "", errors.Wrap(err, "error set redis")
	}
	return s, nil
}

func (r *redisRepository) Delete(key string) error {
	return errors.Wrap(r.cli.Del(r.ctx, key).Err(), "error delete redis")
}

func (r *redisRepository) Expire(key string, expiration time.Duration) error {
	// キーの有効期限を延長
	err := r.cli.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "can't extension ttl")
	}
	return nil

}
