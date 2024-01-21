package handlers

import (
	"log"
)

func (handler *Handler) UpdateEventNotfiyFor() error {

	handler.Repository.Open()

	errEvent := handler.Repository.UpdateEventNotifyFor(
		handler.ModelEvent.UserId,
		handler.ModelEvent.NotifyFor,
		handler.ModelEvent.Id,
	)

	if errEvent != nil {
		log.Println(errEvent)
		return errEvent
	}

	defer handler.Repository.Close()
	return nil
}
