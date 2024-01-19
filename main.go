package main

import (
	"flag"
	tgbot "reminders_tg_bot/telegram-bot"
	"time"
)

func main() {
	var config_path string
	flag.StringVar(&config_path, "config_path", "config.json", "Path to config.json")
	flag.Parse()

	tgBot := &tgbot.TelegramBot{}
	tgBot.Create(config_path)
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