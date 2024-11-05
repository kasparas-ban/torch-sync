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

	// Retry connection with a backoff delay until timeout
	retryInterval := 2 * time.Second
	for {
		if err := db.PingContext(ctx); err == nil {
			DSN = dsn
			DB = db
			return db, nil
		}

		// Check if the context is done (i.e., timeout reached)
		if ctx.Err() != nil {
			return nil, fmt.Errorf("failed to connect to the database after %v: %w", timeout, ctx.Err())
		}

		// Wait for the retry interval before trying again
		fmt.Println("Retrying connection to the database...")
		time.Sleep(retryInterval)
	}
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func handleListenerError(ev pq.ListenerEventType, err error) {
	if err != nil {
		fmt.Printf("Listener error: %v", err)
	}
}
