package redis_session

import (
	"context"

	"github.com/redis/go-redis/v9"
	"golang.org/x/xerrors"
)

type redisSessionID struct {
	cli *redis.Client
	ctx context.Context
}

type RedisSessionID interface {
	Get(key string) (string, error)
}

func NewRedisClient(cli *redis.Client, ctx context.Context) RedisSessionID {
	return &redisSessionID{
		cli: cli,
		ctx: ctx,
	}
}

func (r *redisSessionID) Get(key string) (string, error) {
	s, err := r.cli.Get(r.ctx, key).Result()
	if err != nil {
		return "", xerrors.Errorf("error set redis: %w", err)
	}
	return s, nil
}

func (r *redisSessionID) GetSession(sessionID string) (string, error) {
	s, err := r.Get("fejia")
	if err != nil {
		return "", xerrors.Errorf("redis can't get : %w", err)
	}
	return s, nil
}
