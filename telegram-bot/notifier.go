package telegrambot

import (
	"fmt"
	"log"
	"reminders_tg_got/models"
	"reminders_tg_got/telegram-bot/handlers"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tgBot *TelegramBot) Notifier() {
	handler := &handlers.Handler{
		Config:     tgBot.Config,
		Repository: tgBot.Repository,
	}
	list, getEventsErr := handler.GetEventsWithoutUserid()
	if getEventsErr != nil {
		log.Println(getEventsErr)
	} else {
		alertsNowList := make([]*models.ModelEvents, 0)
		alertsFor := make([]*models.ModelEvents, 0)
		for _, j := range list {
			modelEvent := &models.ModelEvents{
				Id:        j.Id,
				UserId:    j.UserId,
				StartTime: j.StartTime,
				EventName: j.EventName,
				NotifyFor: j.NotifyFor,
			}
			if j.StartTime == time.Now().Unix() {
				alertsFor = append(alertsFor, modelEvent)
			} else if j.StartTime-int64(j.NotifyFor) == time.Now().Unix() {
				alertsNowList = append(alertsNowList, modelEvent)
			}
		}

		var wg sync.WaitGroup

		go func() {
			for _, j := range alertsFor {
				wg.Add(1)
				go func(event *models.ModelEvents) {
					defer wg.Done()
					msg := tgbotapi.NewMessage(event.UserId, fmt.Sprintf(
						"📝: %s\n⏰: %s\nMinutes before the event: %d",
						event.EventName,
						time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
						event.NotifyFor/60,
					))
					tgBot.Bot.Send(msg)
				}(j)
			}
		}()

		go func() {
			for _, j := range alertsNowList {
				wg.Add(1)
				go func(event *models.ModelEvents) {
					defer wg.Done()
					msg := tgbotapi.NewMessage(event.UserId, fmt.Sprintf(
						"%d minutes until this event\n📝: %s\n⏰: %s\n",
						event.NotifyFor/60,
						event.EventName,
						time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
					))
					tgBot.Bot.Send(msg)
				}(j)
			}
		}()

		wg.Wait()
	}
}
