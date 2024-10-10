package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env              string           `yaml:"env" env-default:"local"`
	Storage          storageConfig    `yaml:"storage"`
	JWT              jwtToken         `yaml:"jwt_auth"`
	AccessPermission accessPermission `yaml:"access_permission"`
}

type storageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type jwtToken struct {
	AccessTokenExpiry  string `yaml:"access_token_expiry"`
	RefreshTokenExpiry string `yaml:"refresh_token_expiry"`
}

type accessPermission struct {
	AdminToken string `yaml:"admin_token"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "config.yml", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
