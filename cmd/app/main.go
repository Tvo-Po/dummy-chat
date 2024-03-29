package main

import (
	"context"
	"dummy-chat/internal/manager"
	"dummy-chat/internal/server"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type AppConfig struct {
	Port         string
	LogLevel     slog.Level
	ShutdownTime time.Duration
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
	configData.ShutdownTime *= time.Second
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
	s := server.New(logger, m, config.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go m.Run()
	go s.Serve()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTime)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = errors.Join(err, m.Shutdown(shutdownCtx))
	}()

	go func() {
		defer wg.Done()
		err = errors.Join(err, s.Shutdown(shutdownCtx))
	}()

	wg.Wait()

	if err != nil {
		log.Fatal(err)
	}
}
