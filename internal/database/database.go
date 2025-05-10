package database

import (
	"context"
	"fmt"

	"github.com/Census-Population-Project/API/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBase struct {
	DBPool *pgxpool.Pool
}

func CreateConnectURI(cfg *config.Config) string {
	sslMode := "disable"
	if cfg.Database.SSLMode == true {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		sslMode,
		cfg.Database.MaxConnections,
	)
}

func NewDataBaseClient(cfg *config.Config) (*DataBase, error) {
	pgxPool, err := pgxpool.New(context.Background(), CreateConnectURI(cfg))
	if err != nil {
		return nil, err
	}
	return &DataBase{
		DBPool: pgxPool,
	}, nil
}
