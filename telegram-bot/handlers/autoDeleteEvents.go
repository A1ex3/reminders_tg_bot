package handlers

import "log"

func (handler *Handler) AutoDeleteEvent() error {

	handler.Repository.Open()

	errEvent := handler.Repository.AutoDeleteEvent()
	if errEvent != nil {
		log.Println(errEvent)
		return errEvent
	}

	defer handler.Repository.Close()
	return nil
}
