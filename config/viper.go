package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_URI                string        `mapstructure:"DB_URI"`
	SYMMETRIC_KEY         string        `mapstructure:"SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	PORT                  string        `mapstructure:"PORT"`
	MAILER_PASSWORD       string        `mapstructure:"MAILER_PASSWORD"`
	MAILER_HOST           string        `mapstructure:"MAILER_HOST"`
	GCLOUD_PROJECT_ID     string        `mapstructure:"GCLOUD_PROJECT_ID"`
	GCLOUD_BUCKET_NAME    string        `mapstructure:"GCLOUD_BUCKET_NAME"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	for _, key := range []string{
		"DB_URI", "SYMMETRIC_KEY", "ACCESS_TOKEN_DURATION", "PORT",
		"MAILER_PASSWORD", "MAILER_HOST", "GCLOUD_PROJECT_ID", "GCLOUD_BUCKET_NAME",
	} {
		viper.BindEnv(key)
	}

	// Try reading .env but don't care if it doesn't exist
	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
