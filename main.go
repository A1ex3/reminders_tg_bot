package main

import (
	tgbot "reminders_tg_bot/telegram-bot"
	"time"
)

func main() {

	tgBot := &tgbot.TelegramBot{}
	tgBot.Create("config.json")
	go func() {
		for {
			tgBot.Notifier()
			time.Sleep(time.Second)
		}
	}()
	for {
		tgBot.MessageHandler()
	}
}
