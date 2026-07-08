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
	viper.SetConfigFile(".env")

	// Viper's Unmarshal only honors AutomaticEnv() for keys that are
	// explicitly bound — so bind each field by hand.
	for _, key := range []string{
		"DB_URI", "SYMMETRIC_KEY", "ACCESS_TOKEN_DURATION", "PORT",
		"MAILER_PASSWORD", "MAILER_HOST", "GCLOUD_PROJECT_ID", "GCLOUD_BUCKET_NAME",
	} {
		viper.BindEnv(key)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err // real parse error — still fail
		}
		// no .env file present — fine in prod, real env vars take over
	}

	if err := viper.Unmarshal(config); err != nil {
		return config, err
	}
	return config, nil
}
