package handlers

import (
	"log"
	"reminders_tg_got/models"
)

func (handler *Handler) GetEvents() ([]models.ModelEvents, error) {

	handler.Repository.Open()

	list, errEvent := handler.Repository.GetEventes(
		handler.ModelEvent.UserId,
	)
	if errEvent != nil {
		log.Println(errEvent)
		return nil, errEvent
	}

	defer handler.Repository.Close()
	return list, nil
}

func (handler *Handler) GetEventsWithoutUserid() ([]models.ModelEvents, error) {

	handler.Repository.Open()

	list, errEvent := handler.Repository.GetEventsWithoutUserId()
	if errEvent != nil {
		log.Println(errEvent)
		return nil, errEvent
	}

	defer handler.Repository.Close()
	return list, nil
}