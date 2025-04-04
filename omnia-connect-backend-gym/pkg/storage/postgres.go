package storage

import (
	"api-gateway/internal/config"
	"api-gateway/pkg/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

func PostgresConnect(cfg config.PostgresConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	parseConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Error parsing database config", slog.String("error", err.Error()))
		return nil, err
	}

	parseConfig.MaxConns = 10
	parseConfig.MaxConnIdleTime = 30 * time.Minute

	err = utils.DoWithTries(func() error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pool, err = pgxpool.NewWithConfig(ctx, parseConfig)
		if err != nil {
			slog.ErrorContext(ctx, "Error connecting to database", slog.String("error", err.Error()))
			return err
		}

		if err = pool.Ping(ctx); err != nil {
			slog.Error("Error pinging database", slog.String("error", err.Error()))
			pool.Close()
			return err
		}
		return nil

	}, 5, 5*time.Second)

	return pool, nil
}
