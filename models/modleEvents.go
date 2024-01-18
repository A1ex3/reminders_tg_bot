package models

import (
	"errors"
	"fmt"
	"log"
	"reminders_tg_bot/config"
	errorscustom "reminders_tg_bot/errorsCustom"
	"strconv"
	"strings"
	"time"
)

const (
	eventNameLength int = 128
)

type ModelEvents struct {
	NotifyFor int
	StartTime int64
	UserId    int64
	Id        int64
	EventName string
}

type ModelEventsWithConfig struct {
	ModelEvents *ModelEvents
	Config      *config.Config
}

/*
Gets from the dateTimeStr(sender's local time) parameter a value that contains the date time,

all templates for formatting are stored in the file, config.json dateTimeFormats(array),

offset is the time zone offset from GMT transfer need seconds.

Example:

	var dateTime string = "2024.01.01 12:00:01"
	var offset int64 = 3600
	result, err := parseUnixTime(dateTime string, offset int64) // result = seconds.

Returns unix time(int64)
*/
func (mewc *ModelEventsWithConfig) ParseUnixTime(dateTimeStr string, offset int64) (int64, error) {
	formats := mewc.Config.DateTimeFormats

	var t time.Time
	var err error
	for _, layout := range formats {
		t, err = time.Parse(layout, dateTimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Printf("time parsing error: %v\n", err)
		return 0, errors.New("error: parsing time")
	}

	unixTime := t.Unix() - offset
	return unixTime, nil
}

/*
Gets the seconds from the string that contains the time value.

Supported time units: s - seconds, m - minutes, h - hourses, d - days.

Example:

	var input string = "1h"
	result, err := parseSeconds(input) // result contains 3600 seconds
*/
func (*ModelEventsWithConfig) ParseSeconds(input string) (int, error) {
	var seconds int

	var inputLowerCase string = strings.ToLower(input)
	if strings.Contains(inputLowerCase, "d") {
		days, _ := strconv.Atoi(strings.TrimSuffix(inputLowerCase, "d"))
		seconds += days * 86400
	} else if strings.Contains(inputLowerCase, "m") {
		minutes, _ := strconv.Atoi(strings.TrimSuffix(inputLowerCase, "m"))
		seconds += minutes * 60
	} else if strings.Contains(inputLowerCase, "h") {
		hours, _ := strconv.Atoi(strings.TrimSuffix(inputLowerCase, "h"))
		seconds += hours * 3600
	} else if strings.Contains(inputLowerCase, "s") {
		sec, _ := strconv.Atoi(strings.TrimSuffix(inputLowerCase, "s"))
		seconds += sec
	} else {
		parts := strings.Split(input, ":")
		if len(parts) == 3 {
			hours, _ := strconv.Atoi(parts[0])
			minutes, _ := strconv.Atoi(parts[1])
			sec, _ := strconv.Atoi(parts[2])
			seconds = hours*3600 + minutes*60 + sec
		} else {
			return 0, errors.New(errorscustom.UnsupportedTimeFormat)
		}
	}

	return seconds, nil
}

func (mewc *ModelEventsWithConfig) SetId(value int64) error {
	mewc.ModelEvents.Id = value
	return nil
}

func (mewc *ModelEventsWithConfig) SetUserId(value int64) error {
	mewc.ModelEvents.UserId = value
	return nil
}

func (mewc *ModelEventsWithConfig) SetEventName(value string) error {
	if  len(strings.ReplaceAll(value, " ", "")) < 1 || len(value) > eventNameLength {
		return fmt.Errorf("%s %d characters, length: %d", errorscustom.MaximumLengthEventName, eventNameLength, len(value))
	}
	mewc.ModelEvents.EventName = value
	return nil
}

func (mewc *ModelEventsWithConfig) SetStartTime(startTime string, offset int64) error {
	timeParse, timeParseError := mewc.ParseUnixTime(startTime, offset)

	if timeParseError != nil {
		return timeParseError
	}
	mewc.ModelEvents.StartTime = timeParse
	return nil
}

func (mewc *ModelEventsWithConfig) SetNotifyFor(startTime int64, notifyFor string) error {
	notifyForParse, notifyForParseError := mewc.ParseSeconds(notifyFor)

	if notifyForParseError != nil {
		return notifyForParseError
	} else if int64(notifyForParse)+time.Now().Unix() > startTime {
		log.Println(int64(notifyForParse)+time.Now().Unix(), startTime)
		return errors.New(errorscustom.IncorrectNotificationTime)
	}
	mewc.ModelEvents.NotifyFor = notifyForParse
	return nil
}

func (mewc *ModelEventsWithConfig) Extract(
	notifyFor string,
	startTime string,
	offset int64,
	eventName string,
	userId int64,
) error {
	if err := mewc.SetUserId(userId); err != nil {
		return err
	}
	if err := mewc.SetEventName(eventName); err != nil {
		return err
	}
	if err := mewc.SetStartTime(startTime, offset); err != nil {
		return err
	}
	if err := mewc.SetNotifyFor(mewc.ModelEvents.StartTime, notifyFor); err != nil {
		return err
	}
	return nil
}
