package telegrambot

import (
	config "reminders_tg_bot/config"
	"reminders_tg_bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot        *tgbotapi.BotAPI
	Config     *config.Config
	Repository *database.Repository
}

func (tgBot *TelegramBot) Create(pathToConfig string) {
	config := &config.Config{}
	configErr := config.Unmarshal(pathToConfig)
	if configErr != nil {
		panic(configErr)
	}

	repo := &database.Repository{}
	repo.Config = config

	repo.Open()

	tgBot.Config = config
	tgBot.Repository = repo

	bot, err := tgbotapi.NewBotAPI(tgBot.Config.TgBotApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = tgBot.Config.TgBotDebug
	tgBot.Bot = bot
}
