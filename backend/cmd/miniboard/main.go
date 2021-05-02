package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vrischmann/envconfig"
	yaml "gopkg.in/yaml.v3"

	"github.com/ngalaiko/miniboard/backend"
	"github.com/ngalaiko/miniboard/backend/logger"
)

func main() {
	configPath := flag.String("config", "", "Path to the configuration file, required")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	log := logger.New(logger.Info)
	if *verbose {
		log = logger.New(logger.Debug)
	}

	cfg, err := parseConfiguration(configPath)
	if err != nil {
		log.Fatal("failed to parse configuration: %s", err)
	}

	log.Info("application is starting")

	srv, err := backend.New(log, cfg)
	if err != nil {
		log.Fatal("failed to initialize server: %s", err)
	}

	// Wait for shut down in a separate goroutine.
	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal, 1)
		signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
		sig := <-shutdownCh

		log.Info("received %s, shutting down", sig)

		shutdownTimeout := 15 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		errCh <- srv.Shutdown(shutdownCtx)
	}()

	mainCtx := context.Background()
	if err := srv.Start(mainCtx); err != nil {
		log.Fatal("failed to start the server: %s", err)
	}

	// Handle shutdown errors.
	if err := <-errCh; err != nil {
		log.Error("error during shutdown: %s", err)
	}

	log.Info("application stopped")
}

func parseConfiguration(path *string) (*backend.Config, error) {
	cfg := &backend.Config{}
	if path != nil && *path != "" {
		var err error
		cfg, err = parseConfigurationFromYaml(*path)
		if err != nil {
			return nil, err
		}
	}

	if err := parseConfigurationFromEnvironment(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfigurationFromEnvironment(cfg *backend.Config) error {
	if err := envconfig.InitWithOptions(cfg, envconfig.Options{
		Prefix:      "MINIBOARD",
		AllOptional: true,
	}); err != nil {
		return fmt.Errorf("failed to parse config from env: %w", err)
	}
	return nil
}

func parseConfigurationFromYaml(path string) (*backend.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	cfg := &backend.Config{}
	d := yaml.NewDecoder(file)
	d.KnownFields(true)
	if err := d.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return cfg, nil
}
