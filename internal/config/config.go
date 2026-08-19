package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	PORT             string `mapstructure:"PORT"`
	DATABASE_URL     string `mapstructure:"DATABASE_URL"`
	JWT_SECRET       string `mapstructure:"JWT_SECRET"`
	JWT_EXPIRY_HOURS int    `mapstructure:"JWT_EXPIRY_HOURS"`
	GIN_MODE         string `mapstructure:"GIN_MODE"`
}

// validate the config
func (c Config) Validate() error {
	if c.PORT == "" {
		return errors.New("PORT is required")
	}
	if c.DATABASE_URL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.JWT_SECRET == "" {
		return errors.New("JWT_SECRET is required")
	}
	if c.JWT_EXPIRY_HOURS <= 0 {
		return errors.New(
			"JWT_EXPIRY_HOURS must be greater than zero",
		)
	}
	if c.GIN_MODE != "debug" && c.GIN_MODE != "release" && c.GIN_MODE != "test" {
		return errors.New("GIN_MODE must be debug, release, or test")
	}

	return nil
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("DATABASE_URL")
	_ = viper.BindEnv("JWT_SECRET")
	_ = viper.BindEnv("JWT_EXPIRY_HOURS")
	_ = viper.BindEnv("GIN_MODE")

	cfg := Config{
		PORT:             viper.GetString("PORT"),
		DATABASE_URL:     viper.GetString("DATABASE_URL"),
		JWT_SECRET:       viper.GetString("JWT_SECRET"),
		JWT_EXPIRY_HOURS: viper.GetInt("JWT_EXPIRY_HOURS"),
		GIN_MODE:         viper.GetString("GIN_MODE"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
