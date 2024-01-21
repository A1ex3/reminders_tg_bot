package handlers

import (
	"log"
)

func (handler *Handler) UpdateEventName() error {

	handler.Repository.Open()

	errEvent := handler.Repository.UpdateEventName(
		handler.ModelEvent.UserId,
		handler.ModelEvent.EventName,
		handler.ModelEvent.Id,
	)

	if errEvent != nil {
		log.Println(errEvent)
		return errEvent
	}

	defer handler.Repository.Close()
	return nil
}
