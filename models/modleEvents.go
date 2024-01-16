package models

import (
	"errors"
	"fmt"
	errorscustom "reminders_tg_got/errorsCustom"
	"strconv"
	"strings"
	"time"
)

type ModelEvents struct {
	NotifyFor int
	StartTime int64
	UserId    int64
	Id        int64
	EventName string
}

func (*ModelEvents) parseAndPrintUnixTime(dateTimeStr string, offset int64) (int64, error) {
	formats := []string{
		"02.01.2006 15:04:05",
		"02-01-2006 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"Mon Jan 2 15:04:05 -0700 MST 2006",
		"02 Jan 2006, 15:04:05",
		"2006-01-02 15:04:05",
		"2006.01.02 15:04:05",
	}

	var t time.Time
	var err error
	for _, layout := range formats {
		t, err = time.Parse(layout, dateTimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Printf("time parsing error: %v\n", err)
		return 0, errors.New("error: parsing time")
	}

	unixTime := t.Unix() - offset
	return unixTime, nil
}

func (model *ModelEvents) parseSeconds(input string) (int, error) {
	var seconds int

	if strings.Contains(input, "M") {
		minutes, _ := strconv.Atoi(strings.TrimSuffix(input, "M"))
		seconds += minutes * 60
	} else if strings.Contains(input, "H") {
		hours, _ := strconv.Atoi(strings.TrimSuffix(input, "H"))
		seconds += hours * 3600
	} else if strings.Contains(input, "S") {
		sec, _ := strconv.Atoi(strings.TrimSuffix(input, "S"))
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

func (model *ModelEvents) Extract(
	notifyFor string,
	startTime string,
	offset int64,
	eventName string,
	userId int64,
) error {
	timeParse, timeParseError := model.parseAndPrintUnixTime(startTime, offset)
	notifyForParse, notifyForParseError := model.parseSeconds(notifyFor)

	if timeParseError != nil {
		return timeParseError
	} else if notifyForParseError != nil {
		return notifyForParseError
	} else if int64(notifyForParse)+time.Now().Unix() > timeParse {
		return errors.New(errorscustom.IncorrectNotificationTime)
	} else if len(eventName) > 64 {
		return fmt.Errorf("%s 64 characters", errorscustom.MaximumLengthEventName)
	}

	model.UserId = userId
	model.EventName = eventName
	model.NotifyFor = notifyForParse
	model.StartTime = timeParse
	return nil
}
