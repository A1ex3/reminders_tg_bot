package handlers

import "log"

func (handler *Handler) Registration() bool{
	if !handler.Config.RegistrationAccess{
		return false
	}
	handler.Repository.Open()
	err := handler.Repository.Registration(handler.ModelUser.UserId)
	if err != nil {
		log.Println(err)
		return false
	}
	defer handler.Repository.Close()
	return true
}
