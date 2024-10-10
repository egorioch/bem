package sl

import (
	"fmt"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Logging struct {
		Level  string `yaml:"level"`
		File   string `yaml:"file"`
		Format string `yaml:"format"`
	} `yaml:"log"`
}

type Logger struct {
	slog.Logger
}

func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decoded configFile: %w", err)
	}

	return config, nil
}

func NewLogger(configPath string) (*slog.Logger, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(config.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	jsonHandler := slog.NewJSONHandler(logFile, nil)
	logger := slog.New(jsonHandler)

	return logger, nil
}
