package db

import (
	"context"
	"opm/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

func InitPool(dbURL string) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.MainLogger.Fatalf("❌ Unable to parse database URL: %v\n", err)
	}

	config.MaxConns = 128
	config.MaxConnIdleTime = 10 * time.Minute
	config.ConnConfig.ConnectTimeout = 30 * time.Second

	Conn, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.MainLogger.Fatalf("❌ Unable to create connection pool: %v\n", err)
	}

	logger.MainLogger.Println("✅ Database connection pool initialized")

	// Set timezone
	_, err = Conn.Exec(context.Background(), "SET timezone = 'UTC';")
	if err != nil {
		logger.MainLogger.Fatalf("❌ Error setting timezone: %v\n", err)
	}

	logger.MainLogger.Println("🕒 Timezone set to UTC")
}

// Close closes the database connection pool
func Close() {
	if Conn != nil {
		Conn.Close()
		logger.MainLogger.Println("✅ Database connection pool closed")
	}
}

// Ping checks if the database connection is alive
func Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := Conn.Ping(ctx)
	if err != nil {
		logger.MainLogger.Printf("❌ Database ping failed: %v", err)
		return err
	}
	return nil
}

// BeginTx starts a new transaction
func BeginTx(ctx context.Context) (pgx.Tx, error) {
	return Conn.Begin(ctx)
}

// Query is a convenience wrapper for Conn.Query
func Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return Conn.Query(ctx, sql, args...)
}

// QueryRow is a convenience wrapper for Conn.QueryRow
func QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return Conn.QueryRow(ctx, sql, args...)
}

// Exec is a convenience wrapper for Conn.Exec
func Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return Conn.Exec(ctx, sql, args...)
}
