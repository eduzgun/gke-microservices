package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool   *pgxpool.Pool
	logger *slog.Logger
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(logger *slog.Logger) (*PostgresDB, error) {
	cfg := loadConfig()

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var pool *pgxpool.Pool
	var err error

	// Retry logic
	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pool, err = pgxpool.New(context.Background(), connString)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				logger.Info(
					"connected to postgres",
					"host", cfg.Host,
					"port", cfg.Port,
					"dbname", cfg.DBName,
					"attempt", attempt,
				)
				return &PostgresDB{Pool: pool, logger: logger}, nil
			}
		}

		logger.Error(
			"failed to connect to postgres",
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"host", cfg.Host,
			"port", cfg.Port,
			"error", err,
		)

		if attempt < maxAttempts {
			time.Sleep(5 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}

// Gracefully close connection pool
func (p *PostgresDB) Close() {
	p.logger.Info("closing postgres connection pool")
	p.Pool.Close()
}

func loadConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_NAME"),
		SSLMode:  os.Getenv("PSQL_SSL_MODE"),
	}
}
