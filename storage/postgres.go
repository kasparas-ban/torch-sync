package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var DSN string
var DB *sql.DB

func InitDB(dsn string, timeout time.Duration) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	DSN = dsn
	DB = db
	return db, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func handleListenerError(ev pq.ListenerEventType, err error) {
	if err != nil {
		fmt.Printf("Listener error: %v", err)
	}
}
