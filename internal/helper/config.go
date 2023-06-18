package helper

import (
	"os"

	"github.com/spf13/viper"
)

type JobConfig struct {
	Cron    string `json:"cron"`
	Disable bool   `json:"disable"`
}

type Config struct {
	// r2
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	EndPoint  string `json:"endPoint"`
	Bucket    string `json:"bucket"`
	Prefix    string `json:"prefix"` //  r2 prefix, also used for screenshot file folder path
	Port      string `json:"port"`
	// log
	LogPath  string `json:"logPath"`
	LogLevel string `json:"logLevel"`
	// server
	ApiKey           string `json:"apiKey"`
	EnableMetrics    bool   `json:"enableMetrics"`
	RmImgAfterUpload bool   `json:"rmImgAfterUpload"`
	// scheduler
	ClearLocalImgJob JobConfig `json:"clearLocalImgJob"`
	UpdateR2ImgJob   JobConfig `json:"updateR2ImgJob"`
	// return url
	ReturnUrl string `json:"returnUrl"`
	Headless  bool   `json:"headless"`
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
		ClearLocalImgJob: JobConfig{
			Disable: false,
			Cron:    "0 5 * * * *",
		},
		UpdateR2ImgJob: JobConfig{
			Disable: false,
			Cron:    "0 5 * * * *",
		},
		Headless: true,
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.WatchConfig()

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
