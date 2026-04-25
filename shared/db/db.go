package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func DatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func Open(ctx context.Context) (*sql.DB, error) {
	dsn := DatabaseURL()
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}
