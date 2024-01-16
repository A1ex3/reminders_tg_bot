package database

import (
	"context"
	"database/sql"
)

func (repo *Repository) CreateEvent(
	user_id int64,
	eventName string,
	startTime int64,
	notifyFor int,
) error {
	ctx := context.Background()
	resultCh := make(chan sql.Result)
	errorCh := make(chan error)

	go func() {
		result, err := repo.SqLite.ExecContext(
			ctx,
			"INSERT INTO events VALUES (NULL, $1, $2, $3, $4)",
			user_id,
			eventName,
			startTime,
			notifyFor,
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
