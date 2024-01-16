package telegrambot

import (
	"log"
	"reminders_tg_got/models"
	"reminders_tg_got/telegram-bot/handlers"
	"reminders_tg_got/telegram-bot/tgbutton"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tgBot *TelegramBot) MessageHandler() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tgBot.Bot.GetUpdatesChan(u)

	const (
		InputEventNameState string = "Enter a name for the event"
		InputDateTimeState  string = "Enter the date and time of the event"
		InputNotifyForState string = "Enter how long to give notice"
		NotifyReadyForSend  string = "notifyReadyForSend"
	)
	var (
		startTime string = ""
		eventName string = ""
		eventId   string = ""
	)
	var state string = ""
	log.Println(eventId)
	for update := range updates {
		if update.Message != nil {
			handler := &handlers.Handler{
				Config:     tgBot.Config,
				Repository: tgBot.Repository,
				ModelUser: &models.ModelUser{
					UserId: update.Message.From.ID,
				},
			}
			if update.Message.Text == "/start" {
				handler.Registration()
				state = ""
			} else if update.Message.Text == "/menu" {
				if handler.UserExists() {
					state = ""
					startTime = ""
					eventName = ""
					eventId = ""
					tgButtonInit := tgbutton.Init()
					var tgButton tgbutton.ITgButton = &tgButtonInit
					tgButton.Add(tgButton.Create("New event", "newEvent"))
					buildInlineButtons := tgButton.Build()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "MENU")
					msg.ReplyMarkup = buildInlineButtons
					tgBot.Bot.Send(msg)
				}
			} else if update.Message.Text == "/get" {
				modelEvent := models.ModelEvents{UserId: update.Message.From.ID}
				handler.ModelEvent = &modelEvent
				tgButtonInit := tgbutton.Init()
				var tgButton tgbutton.ITgButton = &tgButtonInit
				list, errEvent := handler.GetEvents()
				if errEvent != nil {
				} else {
					var sb strings.Builder
					for _, j := range list {
						sb.WriteString("EventId:")
						sb.WriteString(strconv.FormatInt(j.Id, 10))
						tgButton.Add(tgButton.Create(j.EventName, sb.String()))
						sb.Reset()
					}
					buildInlineButtons := tgButton.Build()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Events:")
					msg.ReplyMarkup = buildInlineButtons
					tgBot.Bot.Send(msg)
				}
			} else if state == InputDateTimeState {
				if handler.UserExists() {
					msg := tgbotapi.NewMessage(update.Message.From.ID, InputNotifyForState)
					tgBot.Bot.Send(msg)
					startTime = update.Message.Text
					state = InputNotifyForState
				}
			} else if state == InputNotifyForState {
				if handler.UserExists() {
					_, offset := update.Message.Time().Zone()

					modelEvent := models.ModelEvents{}
					extractErr := modelEvent.Extract(
						update.Message.Text,
						startTime,
						int64(offset),
						eventName,
						update.Message.From.ID,
					)
					handler.ModelEvent = &modelEvent
					if extractErr != nil {
						msg := tgbotapi.NewMessage(update.Message.From.ID, extractErr.Error())
						tgBot.Bot.Send(msg)
					} else {
						createEventErr := handler.CreateEvent()
						if createEventErr != nil {
							msg := tgbotapi.NewMessage(update.Message.From.ID, createEventErr.Error())
							tgBot.Bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.From.ID, "Event successfully added")
							tgBot.Bot.Send(msg)
						}
					}
				}
				state = ""
				startTime = ""
				eventName = ""
				eventId = ""
			} else if state == InputEventNameState {
				if update.Message.Sticker != nil {
					state = ""
					msg := tgbotapi.NewMessage(
						update.Message.From.ID,
						"incorrect data",
					)
					tgBot.Bot.Send(msg)
				} else if handler.UserExists() {
					msg := tgbotapi.NewMessage(update.Message.From.ID, InputDateTimeState)
					tgBot.Bot.Send(msg)
					eventName = update.Message.Text
					state = InputDateTimeState
				}
			}
		} else if update.CallbackQuery != nil {
			log.Println(update.CallbackQuery.Data)
			handler := &handlers.Handler{
				Config:     tgBot.Config,
				Repository: tgBot.Repository,
				ModelUser: &models.ModelUser{
					UserId: update.CallbackQuery.From.ID,
				},
			}
			if handler.UserExists() {
				if update.CallbackQuery.Data == "newEvent" {
					startTime = ""
					eventName = ""
					msg := tgbotapi.NewEditMessageText(
						update.CallbackQuery.From.ID,
						update.CallbackQuery.Message.MessageID,
						InputEventNameState,
					)
					state = InputEventNameState
					tgBot.Bot.Send(msg)
				} else if update.CallbackQuery.Data == "deleteEvent" {
					i, err := strconv.ParseInt(eventId, 10, 64)
					if err != nil {
					} else {
						modelEvent := models.ModelEvents{
							UserId: update.CallbackQuery.From.ID,
							Id:     i,
						}
						handler.ModelEvent = &modelEvent
						deleteErr := handler.DeleteEvent()
						if deleteErr != nil {
							msg := tgbotapi.NewMessage(
								update.CallbackQuery.From.ID,
								"failed to delete the event",
							)
							tgBot.Bot.Send(msg)
							eventId = ""
						} else {
							modelEvent := models.ModelEvents{UserId: update.CallbackQuery.From.ID}
							handler.ModelEvent = &modelEvent
							tgButtonInit := tgbutton.Init()
							var tgButton tgbutton.ITgButton = &tgButtonInit
							list, errEvent := handler.GetEvents()
							if errEvent != nil {
							} else {
								var sb strings.Builder
								for _, j := range list {
									sb.WriteString("EventId:")
									sb.WriteString(strconv.FormatInt(j.Id, 10))
									tgButton.Add(tgButton.Create(j.EventName, sb.String()))
									sb.Reset()
								}
								buildInlineButtons := tgButton.Build()
								msg := tgbotapi.NewEditMessageTextAndMarkup(
									update.CallbackQuery.From.ID,
									update.CallbackQuery.Message.MessageID,
									"Events:",
									buildInlineButtons,
								)
								tgBot.Bot.Send(msg)
								eventId = ""
							}
						}
					}
				} else if strings.HasPrefix(update.CallbackQuery.Data, "EventId:") {
					eventId = strings.Split(update.CallbackQuery.Data, "EventId:")[1]
					tgButtonInit := tgbutton.Init()
					var tgButton tgbutton.ITgButton = &tgButtonInit
					tgButton.Add(tgButton.Create("Delete", "deleteEvent"))
					build := tgButton.Build()
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.From.ID,
						update.CallbackQuery.Message.MessageID,
						"Settings:",
						build,
					)
					tgBot.Bot.Send(msg)
				}
			}
		}
	}
}
