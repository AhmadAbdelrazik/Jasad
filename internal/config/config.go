package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Env                   string // Environment Name
	AccessToken           string // Access Token Secret for the JWT
	Dsn                   string // DSN for connection to SQL database
	Port                  string
	LoginAttemptsLimit    int // rate the login attempt per username
	LoginAttemptsDuration int
	RateLimit             int // rate the requests per ip address
	RateDuration          int
}

// NewConfig parses the config from yaml configuration file.
// The default config path are ./configs
func NewConfig(fileName string) (*Configuration, error) {
	var config *Configuration

	viper.AddConfigPath("./configs")
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
