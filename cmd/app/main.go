package main

import (
	"dummy-chat/internal/manager"
	"dummy-chat/internal/server"
	"encoding/json"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
)

type AppConfig struct {
	Port           string `json:"port"`
	LogLevel       slog.Level
}

func ParseAppConfig() (*AppConfig, error) {
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	configData := &AppConfig{}
	err = json.Unmarshal(configFile, configData)
	if err != nil {
		return nil, err
	}
	return configData, nil
}

func main() {
	config, err := ParseAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	opts := &slog.HandlerOptions{AddSource: true, Level: config.LogLevel}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	m := manager.New(logger)
	serv := server.New(logger, m, config.Port)

	go m.Run()
	serv.Serve()
}
