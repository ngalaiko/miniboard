package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

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

	mainCtx := context.Background()

	_, err = server.New(mainCtx, logger, cfg)
	if err != nil {
		logger.Fatal("failed to initialize server: %s", err)
	}
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
