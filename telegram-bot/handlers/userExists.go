package handlers

import (
	"log"
)

func (handler *Handler) UserExists() bool {
	handler.Repository.Open()
	err := handler.Repository.UserExists(handler.ModelUser.UserId)
	if err != nil {
		log.Println(err)
		return false
	}
	
	defer handler.Repository.Close()
	return true
}