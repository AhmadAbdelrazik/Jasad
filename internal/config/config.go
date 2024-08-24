package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Env                string
	AccessToken        string
	Dsn                string
	Port               string
	LoginAttemptsLimit int
}

func NewConfig() (*Configuration, error) {
	var config *Configuration

	viper.AddConfigPath("./configs")
	viper.SetConfigName("development")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
