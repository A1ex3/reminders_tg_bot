package database

import (
	"context"
	"database/sql"
)

func (repo *Repository) DeleteEvent(
	userId int64,
	eventID int64,
) error {
	ctx := context.Background()
	resultCh := make(chan sql.Result)
	errorCh := make(chan error)

	go func() {
		result, err := repo.SqLite.ExecContext(
			ctx,
			"DELETE FROM events WHERE e_user_id = $1 AND id = $2",
			userId,
			eventID,
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
