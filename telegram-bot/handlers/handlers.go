package handlers

import (
	"reminders_tg_got/config"
	"reminders_tg_got/database"
	"reminders_tg_got/models"
)

type Handler struct {
	Config     *config.Config
	Repository *database.Repository
	ModelUser *models.ModelUser
	ModelEvent *models.ModelEvents
}
