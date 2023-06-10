package helper

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AccessKey     string `json:"accessKey"`
	SecretKey     string `json:"secretKey"`
	EndPoint      string `json:"endPoint"`
	Bucket        string `json:"bucket"`
	Prefix        string `json:"prefix"` //  r2 prefix, also used for screenshot file folder path
	Port          string `json:"port"`
	ApiKey        string `json:"apiKey"`
	EnableMetrics bool   `json:"enableMetrics"`
	LogPath       string `json:"logPath"`
	LogLevel      string `json:"logLevel"`
}

func initDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func SetupConfig() (*Config, error) {
	// default config
	cfg := Config{
		Prefix:        "screenshot",
		Port:          "8080",
		EnableMetrics: false,
		LogPath:       "log",
		LogLevel:      "info",
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	initDir(cfg.LogPath)
	initDir(cfg.Prefix)

	return &cfg, nil
}
