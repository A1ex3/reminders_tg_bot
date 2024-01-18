package telegrambot

import (
	"fmt"
	"log"
	"reminders_tg_bot/models"
	"reminders_tg_bot/telegram-bot/handlers"
	"reminders_tg_bot/telegram-bot/tgbutton"
	"strconv"
	"strings"
	"time"

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
		eventId string = ""
	)
	var state string = "" // A variable that is responsible for the state of the user's actions.

	modelEvent := &models.ModelEventsWithConfig{
		Config:      tgBot.Config,
		ModelEvents: &models.ModelEvents{},
	}

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
				/*
					When a user enters "/start", the "Registration()" method will try to register the user,
					if registration is disabled or the user already exists, nothing will happen.
				*/
				handler.Registration()
				state = ""
			} else if update.Message.Text == "/menu" {
				/*
					Outputs inline buttons with various functions.
				*/
				if handler.UserExists() {
					state = ""
					eventId = ""
					tgButtonInit := tgbutton.Init()
					var tgButton tgbutton.ITgButton = &tgButtonInit
					tgButton.Add(tgButton.Create("New event", "newEvent"))
					buildInlineButtons := tgButton.Build()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu:")
					msg.ReplyMarkup = buildInlineButtons
					tgBot.Bot.Send(msg)
				}
			} else if update.Message.Text == "/get" {
				/*
					Outputs all events that the user has created
				*/
				if handler.UserExists() {
					modelEvent := models.ModelEvents{UserId: update.Message.From.ID}
					handler.ModelEvent = &modelEvent
					tgButtonInit := tgbutton.Init()
					var tgButton tgbutton.ITgButton = &tgButtonInit
					list, errEvent := handler.GetEvents()
					if errEvent != nil {
					} else {
						var sb strings.Builder
						for _, j := range list {
							/*
								A callBackData is bound to each button,
								which contains the value of "EventId:<id of the event>"".
							*/
							sb.WriteString("EventId:")
							sb.WriteString(strconv.FormatInt(j.Id, 10))
							tgButton.Add(tgButton.Create(j.EventName, sb.String()))
							sb.Reset()
						}
						buildInlineButtons := tgButton.Build()
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Events:")
						msg.ReplyMarkup = buildInlineButtons
						tgBot.Bot.Send(msg)
						state = ""
					}
				}
			} else if state == InputDateTimeState {
				if handler.UserExists() {
					_, offset := update.Message.Time().Zone()
					modelEvent.ModelEvents.UserId = update.Message.From.ID
					err := modelEvent.SetStartTime(update.Message.Text, int64(offset))
					if err != nil {
						msg := tgbotapi.NewMessage(
							update.Message.From.ID,
							fmt.Sprintf(
								"⚠️%s\n\n📄List of supported formats:\n%s",
								err.Error(),
								func() string {
									var sb strings.Builder
									for _, j := range tgBot.Config.DateTimeFormats {
										sb.WriteString(j)
										sb.WriteString(";\n")
									}
									return sb.String()
								}(),
							),
						)
						tgBot.Bot.Send(msg)
					} else {
						handler.ModelEvent = modelEvent.ModelEvents
						msg := tgbotapi.NewMessage(update.Message.From.ID, InputNotifyForState)
						tgBot.Bot.Send(msg)
						//log.Println("InputDateTimeState:", modelEvent.ModelEvents, update.Message.Text)
						state = InputNotifyForState
					}
				}
			} else if state == InputNotifyForState {
				if handler.UserExists() {
					modelEvent.ModelEvents.UserId = update.Message.From.ID
					err := modelEvent.SetNotifyFor(modelEvent.ModelEvents.StartTime, update.Message.Text)

					//log.Println("InputNotifyForState:", modelEvent.ModelEvents)

					handler.ModelEvent = modelEvent.ModelEvents
					if err != nil {
						msg := tgbotapi.NewMessage(
							update.Message.From.ID,
							fmt.Sprintf(
								"⚠️%s\nAllowed formats:\ns - seconds;\nm - minutes;\nh - hours;\nd - days;\nh:m:s;\nExamples:\n2h - will notify 2 hours before the event;",
								err.Error(),
							),
						)
						tgBot.Bot.Send(msg)
					} else {
						createEventErr := handler.CreateEvent()
						if createEventErr != nil {
							msg := tgbotapi.NewMessage(update.Message.From.ID, createEventErr.Error())
							tgBot.Bot.Send(msg)
							state = ""
							eventId = ""
						} else {
							msg := tgbotapi.NewMessage(update.Message.From.ID, "Event successfully added")
							tgBot.Bot.Send(msg)
							state = ""
							eventId = ""
						}
					}
				}
			} else if state == InputEventNameState {
				if handler.UserExists() {
					if update.Message.Sticker != nil {
						msg := tgbotapi.NewMessage(
							update.Message.From.ID,
							"incorrect data",
						)
						tgBot.Bot.Send(msg)
					} else {
						modelEvent.ModelEvents.UserId = update.Message.From.ID
						err := modelEvent.SetEventName(update.Message.Text)
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.From.ID, err.Error())
							tgBot.Bot.Send(msg)
						} else {
							handler.ModelEvent = modelEvent.ModelEvents
							msg := tgbotapi.NewMessage(update.Message.From.ID, InputDateTimeState)
							tgBot.Bot.Send(msg)

							//log.Println("InputEventNameState:", modelEvent.ModelEvents, update.Message.Text)

							state = InputDateTimeState
						}
					}
				}
			}
		} else if update.CallbackQuery != nil {
			handler := &handlers.Handler{
				Config:     tgBot.Config,
				Repository: tgBot.Repository,
				ModelUser: &models.ModelUser{
					UserId: update.CallbackQuery.From.ID,
				},
			}
			if handler.UserExists() {
				if update.CallbackQuery.Data == "newEvent" {
					/*
						If "newEvent" was sent, then the "Menu:" message is edited
						and a message asking for the event name is displayed,
						then the state gets the InputEventNameState.
					*/
					msg := tgbotapi.NewEditMessageText(
						update.CallbackQuery.From.ID,
						update.CallbackQuery.Message.MessageID,
						InputEventNameState,
					)
					state = InputEventNameState
					tgBot.Bot.Send(msg)
				} else if update.CallbackQuery.Data == "backToList" {
					modelEvent := models.ModelEvents{UserId: update.CallbackQuery.From.ID}
					handler.ModelEvent = &modelEvent
					tgButtonInit := tgbutton.Init()
					var tgButton tgbutton.ITgButton = &tgButtonInit
					list, errEvent := handler.GetEvents()
					if errEvent != nil {
					} else {
						var sb strings.Builder
						for _, j := range list {
							/*
								A callBackData is bound to each button,
								which contains the value of "EventId:<id of the event>"".
							*/
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
					}
				} else if update.CallbackQuery.Data == "deleteEvent" {
					/*
						This is where the event is deleted. If the deletion was successful,
						the message is edited and the event list is displayed again.
					*/
					i, err := strconv.ParseInt(eventId, 10, 64)
					if err != nil {
						msg := tgbotapi.NewEditMessageText(
							update.CallbackQuery.From.ID,
							update.CallbackQuery.Message.MessageID,
							"Settings:",
						)
						tgBot.Bot.Send(msg)
						eventId = ""
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
					/*
						If the user clicked on an event (inline button), the prefix is checked which should contain EventId:
						and after it the event id, if all was successful, the "Events:" message is edited and displayed
						from the inline buttons. And also the record information is displayed.

							*******************
							| <Event intfo>:  |
							| ********* 	  |
							| |Delete|		  |
							| *********		  |
							| |Back  |        |
							| ********        |
							*******************
					*/
					modelEvent := models.ModelEvents{UserId: update.CallbackQuery.From.ID}
					handler.ModelEvent = &modelEvent

					row, err := handler.GetEvents()

					if err != nil {
						log.Printf(err.Error())
					} else {
						result := make([]models.ModelEvents, 0)

						eventId = strings.Split(update.CallbackQuery.Data, "EventId:")[1]

						for _, j := range row {
							if val := strconv.Itoa(int(j.Id)); val == eventId {
								result = append(result, j)
							}
						}

						fmtResult := fmt.Sprintf(
							"📝: %s\n⏰: %s\n⏲: %s",
							result[0].EventName,
							time.Unix(result[0].StartTime, 0).Format("2006-01-02 15:04:05"),
							time.Unix(result[0].StartTime-int64(result[0].NotifyFor), 0).Format("2006-01-02 15:04:05"),
						)
						tgButtonInit := tgbutton.Init()
						var tgButton tgbutton.ITgButton = &tgButtonInit
						tgButton.Add(tgButton.Create("Delete", "deleteEvent"))
						tgButton.Add(tgButton.Create("Back", "backToList"))
						build := tgButton.Build()
						msg := tgbotapi.NewEditMessageTextAndMarkup(
							update.CallbackQuery.From.ID,
							update.CallbackQuery.Message.MessageID,
							fmtResult,
							build,
						)
						tgBot.Bot.Send(msg)
					}
				}
			}
		}
	}
}
