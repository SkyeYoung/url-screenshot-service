package helper

import (
	"github.com/spf13/viper"
)

type Config struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	EndPoint  string `json:"endPoint"`
	Bucket    string `json:"bucket"`
	Prefix    string `json:"prefix"`
	Port      string `json:"port"`
}

func SetupConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
