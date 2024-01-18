package handlers

import (
	"reminders_tg_bot/config"
	"reminders_tg_bot/database"
	"reminders_tg_bot/models"
)

type Handler struct {
	Config     *config.Config
	Repository *database.Repository
	ModelUser *models.ModelUser
	ModelEvent *models.ModelEvents
}
