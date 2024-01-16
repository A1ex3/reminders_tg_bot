package database

import (
	"context"
	"database/sql"
)

func (repo *Repository) Registration(user_id int64) error {
	ctx := context.Background()
	resultCh := make(chan sql.Result)
	errorCh := make(chan error)

	defer close(resultCh)
	defer close(errorCh)

	go func() {
		result, err := repo.SqLite.ExecContext(
			ctx,
			"INSERT INTO users VALUES ($1)",
			user_id,
		)
		if err != nil {
			errorCh <- err
		} else {
			resultCh <- result
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errorCh:
		return err
	case <-resultCh:
		return nil
	}
}
