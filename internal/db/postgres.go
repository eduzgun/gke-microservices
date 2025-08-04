package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB() (*PostgresDB, error) {
	// Get configuration from environment or use defaults
	dbConfig := DBConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     "5432", // Default PostgreSQL port
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_NAME"),
		SSLMode:  "disable",
	}

	// Build connection string
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode)

	var pool *pgxpool.Pool
	var err error

	// Retry logic
	maxAttempts := 10
	for i := range maxAttempts {
		pool, err = pgxpool.New(context.Background(), connString)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				log.Info().Msg("established connection to postgres db")
				return &PostgresDB{Pool: pool}, nil
			}
		}
		log.Warn().Err(err).Msgf("Failed to connect to database, retrying in 5 seconds... (%d/%d)", i+1, maxAttempts)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}

func (p *PostgresDB) Close() {
	p.Pool.Close()
}

type RedisDB struct {
	Client *redis.Client
}

func NewRedisDB() (*RedisDB, error) {
	// Redis configuration
	redisConfig := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSW"),
		DB:       0, // Use default DB
	}

	var client *redis.Client
	var err error

	// Retry logic
	maxAttempts := 10
	for i := range maxAttempts {
		client = redis.NewClient(redisConfig)

		// Ping Redis to verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = client.Ping(ctx).Result()
		cancel()

		if err == nil {
			log.Info().Msg("established connection to Redis")
			return &RedisDB{Client: client}, nil
		}

		log.Warn().Err(err).Msgf("Failed to connect to Redis, retrying in 5 seconds... (%d/%d)", i+1, maxAttempts)
		time.Sleep(5 * time.Second)

		// Clean up the failed client
		if client != nil {
			_ = client.Close()
		}
	}

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %w", maxAttempts, err)
}

func (r *RedisDB) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}
