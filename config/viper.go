package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	PORT string `mapstructure:"PORT"`
	DB_USERNAME string `mapstructure:"DB_USERNAME"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_HOSTNAME string `mapstructure:"DB_HOSTNAME"`
	DB_PORT string `mapstructure:"DB_PORT"`
	DB_NAME string `mapstructure:"DB_NAME"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			PORT: os.Getenv("PORT"),
			DB_USERNAME: os.Getenv("DB_USERNAME"),
			DB_PASSWORD: os.Getenv("DB_PASSWORD"),
			DB_HOSTNAME: os.Getenv("DB_HOSTNAME"),
			DB_PORT: os.Getenv("DB_PORT"),
			DB_NAME: os.Getenv("DB_NAME"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.PORT == "" {
		err = errors.New("PORT is required")
		return
	}
	
	if config.DB_USERNAME == "" {
		err = errors.New("DB_USERNAME is required")
		return
	}

	if config.DB_PASSWORD == "" {
		err = errors.New("DB_PASSWORD is required")
		return
	}

	if config.DB_HOSTNAME == "" {
		err = errors.New("DB_HOSTNAME is required")
		return
	}

	if config.DB_PORT == "" {
		err = errors.New("DB_PORT is required")
		return
	}

	if config.DB_NAME == "" {
		err = errors.New("DB_NAME is required")
		return
	}

	return
}