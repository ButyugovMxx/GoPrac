package repository

import (
	"context"
	"fmt"
	"os"
)

func Migrate(db *MovieDbConnect, ctx context.Context) error {

	data, err := os.ReadFile("migrations/001_create_movies.sql")
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	query := string(data)
	_, err = db.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}

	return nil
}
