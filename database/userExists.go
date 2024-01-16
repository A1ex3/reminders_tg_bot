package database

import (
	"context"
	"errors"
)

type userExists struct {
	value bool
}

func (repo *Repository) UserExists(user_id int64) error {
	ctx := context.Background()

	userExistsCh := make(chan userExists)
	errorCh := make(chan error)

	go func() {
		defer close(errorCh)
		defer close(userExistsCh)

		row := repo.SqLite.QueryRowContext(
			ctx,
			"SELECT TRUE FROM users WHERE user_id = $1",
			user_id,
		)
		userExists := userExists{}
		err := row.Scan(&userExists.value)
		if err != nil {
			errorCh <- err
		} else {
			userExistsCh <- userExists
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errorCh:
		return err
	case uE, ok := <-userExistsCh:
		if !ok {
			return errors.New("unexpected issue while checking user existence")
		}
		if uE.value {
			return nil
		} else {
			return errors.New("record not found")
		}
	}
}