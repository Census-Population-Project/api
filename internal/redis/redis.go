package redis

import (
	"strconv"

	"github.com/Census-Population-Project/API/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       db,
	})
}
