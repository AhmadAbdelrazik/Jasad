package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Env         string
	AccessToken string
	Dsn         string
	Port        string
}

func NewConfig() (*Configuration, error) {
	var config *Configuration

	// Default Values
	viper.SetDefault("env", "development")
	viper.SetDefault("accessToken", "69zDfhhZUxnNl63VqmV3EQWja9++RsqORbltMyeTMVHm")
	viper.SetDefault("dsn", "ahmad:password@/jasad?parseTime=true")
	viper.SetDefault("port", ":3000")

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
