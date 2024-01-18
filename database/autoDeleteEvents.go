package database

import (
	"context"
	"database/sql"
)

/*
This method automatically deletes records that have been expired for more than 90 days.
*/
func (repo *Repository) AutoDeleteEvent() error {
	ctx := context.Background()
	resultCh := make(chan sql.Result)
	errorCh := make(chan error)

	go func() {
		result, err := repo.SqLite.ExecContext(
			ctx,
			"DELETE FROM events WHERE CAST(strftime('%s', 'now') AS INTEGER) > events.start_time + 7776000", // 7776000 seconds = 90 days.
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
