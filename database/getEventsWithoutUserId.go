package database

import (
	"context"
	"errors"
	"log"
	"reminders_tg_got/models"
)

func (repo *Repository) GetEventsWithoutUserId() ([]models.ModelEvents, error) {
	ctx := context.Background()

	modelEventsCh := make(chan []models.ModelEvents)
	errorCh := make(chan error)

	go func() {
		defer close(errorCh)
		defer close(modelEventsCh)

		rows, err := repo.SqLite.QueryContext(
			ctx,
			"SELECT id, e_user_id, event_name, start_time, notify_for FROM events",
		)

		list := []models.ModelEvents{}
	
		if err != nil {
			errorCh <- err
		} else {
			for rows.Next(){
				me := models.ModelEvents{}
				errS := rows.Scan(&me.Id, &me.UserId, &me.EventName, &me.StartTime, &me.NotifyFor)
				if errS != nil{
					log.Println(err)
					continue
				}
				list = append(list, me)
			}
			modelEventsCh <- list
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errorCh:
		return nil, err
	case uE, ok := <-modelEventsCh:
		if !ok {
			return nil, errors.New("unexpected issue while checking user existence")
		}
		if ok {
			return uE, nil
		} else {
			return nil, errors.New("records not found")
		}
	}
}
