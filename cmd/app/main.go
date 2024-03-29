package main

import (
	"context"
	"dummy-chat/internal/manager"
	"dummy-chat/internal/server"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type AppConfig struct {
	LogLevel     slog.Level    `envconfig:LOG_LEVEL`
	ShutdownTime time.Duration `envconfig:SHUTDOWN_TIME`
}

func ParseAppConfig() (*AppConfig, error) {
	configData := &AppConfig{}
	err := envconfig.Process("", configData)
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
	s := server.New(logger, m, "8000")

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
