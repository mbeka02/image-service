package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_URI string `mapstructure:"DB_URI"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	// tell Viper the location of the config file
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	// read values
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return config, err
	}
	return config, nil
}
