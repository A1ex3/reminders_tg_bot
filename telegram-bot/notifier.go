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

func (*TelegramBot) buildMessage(event *models.ModelEvents) string {

	minutesBefore := event.NotifyFor
	beforeText := "until the event"
	if time.Now().Unix() > event.StartTime-int64(event.NotifyFor) {
		event.NotifyFor = 0
		beforeText = "the event has started"
	}
	if event.NotifyFor < 60 {
		return fmt.Sprintf(
			"🔊: %d seconds %s\n📝: %s\n⏰: %s",
			event.NotifyFor,
			beforeText,
			event.EventName,
			time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
		)
	} else if event.NotifyFor < 3600 {
		return fmt.Sprintf(
			"🔊: %d minutes %s\n📝: %s\n⏰: %s",
			event.NotifyFor/60,
			beforeText,
			event.EventName,
			time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
		)
	} else if event.NotifyFor < 86400 {
		hoursBefore := event.NotifyFor / 3600
		minutesBefore = (event.NotifyFor % 3600) / 60
		return fmt.Sprintf(
			"🔊: %d hours and %d minutes %s\n📝: %s\n⏰: %s",
			hoursBefore,
			minutesBefore,
			beforeText,
			event.EventName,
			time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
		)
	} else if event.NotifyFor < 2_592_000 {
		days := event.NotifyFor / 86400
		weeks := days / 7
		days = days % 7
		return fmt.Sprintf(
			"🔊: %d weeks and %d days %s\n📝: %s\n⏰: %s",
			weeks, days,
			beforeText,
			event.EventName,
			time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
		)
	} else {
		years := event.NotifyFor / 31104000
		weeks := (event.NotifyFor % 31104000) / 604800
		days := ((event.NotifyFor % 31104000) % 604800) / 86400
		return fmt.Sprintf(
			"🔊: %d years, %d weeks and %d days %s\n📝: %s\n⏰: %s",
			years, weeks, days,
			beforeText,
			event.EventName,
			time.Unix(event.StartTime, 0).Format("2006-01-02 15:04:05"),
		)
	}
}

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
					msg := tgbotapi.NewMessage(event.UserId, tgBot.buildMessage(event))
					tgBot.Bot.Send(msg)
				}(j)
			}
		}()

		go func() {
			for _, j := range alertsNowList {
				wg.Add(1)
				go func(event *models.ModelEvents) {
					defer wg.Done()
					msg := tgbotapi.NewMessage(event.UserId, tgBot.buildMessage(event))
					tgBot.Bot.Send(msg)
				}(j)
			}
		}()
		
		wg.Add(1)
		go func(){
			defer wg.Done()
			err := handler.AutoDeleteEvent()
			if err != nil{
				log.Println(err)
			}
		}()

		wg.Wait()
	}
}
