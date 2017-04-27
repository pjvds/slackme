package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Config struct {
	DefaultChannel string
	Channels       []Channel
}

type Channel struct {
	Name       string
	WebhookUrl string
}

func (this Channel) Post(message string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"text": message,
	})
	response, err := http.Post(this.WebhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	return nil
}

func LoadConfig() (Config, error) {
	var config Config
	file, err := os.Open(configPath())

	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func StoreConfig(config Config) error {
	file, err := os.OpenFile(configPath(), os.O_CREATE|os.O_TRUNC|os.O_RDONLY, 0655)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func configPath() string {
	return os.ExpandEnv("$HOME/.report")

}