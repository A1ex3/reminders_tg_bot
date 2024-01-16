package handlers

import "log"

func (handler *Handler) DeleteEvent() error {

	handler.Repository.Open()

	errEvent := handler.Repository.DeleteEvent(
		handler.ModelEvent.UserId,
		handler.ModelEvent.Id,
	)
	if errEvent != nil {
		log.Println(errEvent)
		return errEvent
	}

	defer handler.Repository.Close()
	return nil
}
