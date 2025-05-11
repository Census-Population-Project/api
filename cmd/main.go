package main

import (
	"context"
	"sync"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"
	"github.com/Census-Population-Project/API/internal/logger"
	"github.com/Census-Population-Project/API/internal/redis"
	"github.com/Census-Population-Project/API/internal/service/api"
)

func main() {
	var wg sync.WaitGroup
	log := logger.NewLogger()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Info("Config (env) loaded!")

	rdb := redis.NewRedisClient(cfg, cfg.Redis.DB)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Info("Connected to Redis!")

	db, err := database.NewDataBaseClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	if err := db.DBPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Info("Connected to database!")

	server := api.NewServerHttp(log, cfg, db, rdb, &wg)

	err = server.UsersService.InitDefaultUser()
	if err != nil {
		log.Fatalf("Failed to initialize default user: %v", err)
	}

	server.InitAPI()
	server.Start()

	wg.Wait()
}
