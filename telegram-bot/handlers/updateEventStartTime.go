package handlers

import (
	"log"
)

func (handler *Handler) UpdateEventStartTime() error {

	handler.Repository.Open()

	errEvent := handler.Repository.UpdateEventStartTime(
		handler.ModelEvent.UserId,
		handler.ModelEvent.StartTime,
		handler.ModelEvent.Id,
	)

	if errEvent != nil {
		log.Println(errEvent)
		return errEvent
	}

	defer handler.Repository.Close()
	return nil
}
