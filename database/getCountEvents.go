package database

import (
	"context"
	"errors"
)

type countEvents struct {
	count int
}

func (repo *Repository) GetCountEvents(user_id int64) (int ,error) {
	ctx := context.Background()

	countCh := make(chan countEvents)
	errorCh := make(chan error)

	go func() {
		defer close(errorCh)
		defer close(countCh)

		row := repo.SqLite.QueryRowContext(
			ctx,
			"SELECT count(*) FROM events WHERE e_user_id = $1",
			user_id,
		)
		userExists := countEvents{}
		err := row.Scan(&userExists.count)
		if err != nil {
			errorCh <- err
		} else {
			countCh <- userExists
		}
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err := <-errorCh:
		return 0, err
	case uE, ok := <-countCh:
		if !ok {
			return 0, errors.New("unexpected issue while checking user existence")
		}
		if ok{
			return uE.count, nil
		} else {
			return 0, errors.New("record not found")
		}
	}
}