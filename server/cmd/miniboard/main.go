package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/ngalaiko/miniboard/server"
	"github.com/ngalaiko/miniboard/server/logger"
)

func main() {
	configPath := flag.String("config", "", "Path to the configuration file, required")
	flag.Parse()

	logger := logger.New()

	if *configPath == "" {
		logger.Fatal("--config is not defined")
	}

	cfg, err := parseConfiguration(*configPath)
	if err != nil {
		logger.Fatal("failed to parse configuration: %s", err)
	}

	logger.Info("application is starting")

	srv, err := server.New(logger, cfg)
	if err != nil {
		logger.Fatal("failed to initialize server: %s", err)
	}

	// Wait for shut down in a separate goroutine.
	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal)
		signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
		sig := <-shutdownCh

		logger.Info("received %s, shutting down", sig)

		shutdownTimeout := 15 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		errCh <- srv.Shutdown(shutdownCtx)
	}()

	mainCtx := context.Background()
	if err := srv.Start(mainCtx); err != nil {
		logger.Fatal("failed to start the server: %s", err)
	}

	// Handle shutdown errors.
	if err := <-errCh; err != nil {
		logger.Warn("error during shutdown: %s", err)
	}

	logger.Info("application stopped")
}

func parseConfiguration(path string) (*server.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	cfg := &server.Config{}
	if err := yaml.UnmarshalStrict(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return cfg, nil
}
