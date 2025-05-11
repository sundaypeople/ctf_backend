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

// PushToList は指定したリストキーに要素を追加します。
func (r *redisRepository) SetList(key string, elements ...interface{}) error {
	err := r.cli.LPush(r.ctx, key, elements...).Err()
	if err != nil {
		return fmt.Errorf("LPUSHに失敗しました: %w", err)
	}
	fmt.Printf("リスト '%s' に要素を追加しました\n", key)
	return nil
}

// SetExpiration は指定したキーに有効期限を設定します。
func (r *redisRepository) SetExpiration(key string, expiration time.Duration) error {
	err := r.cli.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("Expireの設定に失敗しました: %w", err)
	}
	fmt.Printf("キー '%s' に有効期限を %v に設定しました\n", key, expiration)
	return nil
}
