package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"
)

func pingWithTimeout(ping func() error, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- ping()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func closeSQLDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func homeDir() (string, error) {
	if h := os.Getenv("HOME"); h != "" {
		return h, nil
	}
	if h := os.Getenv("USERPROFILE"); h != "" {
		return h, nil
	}
	return "", errors.New("home directory not found")
}
