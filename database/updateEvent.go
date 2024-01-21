package database

import (
	"context"
	"database/sql"
)

func (repo *Repository) UpdateEventName(
	userId int64,
	eventName string,
	eventId int64,
) error {
	ctx := context.Background()
	resultCh := make(chan sql.Result)
	errorCh := make(chan error)

	go func() {
		result, err := repo.SqLite.ExecContext(
			ctx,
			"UPDATE events SET event_name = $1 WHERE id = $2 AND e_user_id = $3",
			eventName,
			eventId,
			userId,
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

func (repo *Repository) UpdateEventStartTime(
	userId int64,
	startTime int64,
	eventId int64,
) error {
	ctx := context.Background()
	resultCh := make(chan sql.Result)
	errorCh := make(chan error)

	go func() {
		result, err := repo.SqLite.ExecContext(
			ctx,
			"UPDATE events SET start_time = $1 WHERE id = $2 AND e_user_id = $3",
			startTime,
			eventId,
			userId,
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