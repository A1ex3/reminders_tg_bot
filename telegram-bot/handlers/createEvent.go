package handlers

import (
	"fmt"
	"log"
)

func (handler *Handler) CreateEvent() error {

	handler.Repository.Open()

	count, errCount := handler.Repository.GetCountEvents(handler.ModelEvent.UserId)

	if errCount != nil {
		return errCount
	}

	if count > handler.Config.MaxCountEventsPerUser {
		return fmt.Errorf("record limit exceeded, maximum %d", handler.Config.MaxCountEventsPerUser)
	}

	errEvent := handler.Repository.CreateEvent(
		handler.ModelEvent.UserId,
		handler.ModelEvent.EventName,
		handler.ModelEvent.StartTime,
		handler.ModelEvent.NotifyFor,
	)
	if errEvent != nil {
		log.Println(errEvent)
		return errEvent
	}

	defer handler.Repository.Close()
	return nil
}
