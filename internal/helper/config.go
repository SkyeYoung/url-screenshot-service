package helper

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	EndPoint  string `json:"endPoint"`
	Bucket    string `json:"bucket"`
	Prefix    string `json:"prefix"` //  r2 prefix, also used for screenshot file folder path
	Port      string `json:"port"`

	LogPath  string `json:"logPath"`
	LogLevel string `json:"logLevel"`

	ApiKey           string `json:"apiKey"`
	EnableMetrics    bool   `json:"enableMetrics"`
	RmImgAfterUpload bool   `json:"rmImgAfterUpload"`
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
		Prefix:           "screenshot",
		Port:             "8080",
		ApiKey:           "123456",
		LogPath:          "log",
		LogLevel:         "info",
		EnableMetrics:    false,
		RmImgAfterUpload: false,
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
