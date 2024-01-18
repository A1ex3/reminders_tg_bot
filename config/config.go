package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	TgBotDebug            bool     `json:"tgBotDebug"`
	RegistrationAccess    bool     `json:"registrationAccess"`
	MaxCountEventsPerUser int      `json:"maxCountEventsPerUser"`
	TgBotApiToken         string   `json:"tgBotApiToken"`
	PathToDataBase        string   `json:"pathToDataBase"`
	DateTimeFormats       []string `json:"dateTimeFormats"`
}

func (conf *Config) Unmarshal(path_to_config string) error {
	file, err := os.Open(path_to_config)

	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	byteValue, readErr := io.ReadAll(file)

	if readErr != nil{
		panic(readErr)
	}

	json.Unmarshal(byteValue, &conf)
	return nil
}
