package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	TgBotDebug            bool   `json:"tgBotDebug"`
	RegistrationAccess    bool   `json:"registrationAccess"`
	MaxCountEventsPerUser int    `json:"maxCountEventsPerUser"`
	TgBotApiToken         string `json:"tgBotApiToken"`
	PathToDataBase        string `json:"pathToDataBase"`
}

func (conf *Config) Unmarshal(path_to_config string) error {
	file, err := os.Open(path_to_config)

	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &conf)
	return nil
}
