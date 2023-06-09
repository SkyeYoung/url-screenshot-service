package helper

import (
	"github.com/spf13/viper"
)

type Config struct {
	AccessKey string
	SecretKey string
	EndPoint  string
	Bucket    string
	Prefix    string
	Port      string
}

func GetConfig() (*Config, error) {
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
