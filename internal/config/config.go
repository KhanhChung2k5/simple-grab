package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PORT             string `mapstructure:"PORT"`
	DATABASE_URL     string `mapstructure:"DATABASE_URL"`
	JWT_SECRET       string `mapstructure:"JWT_SECRET"`
	JWT_EXPIRY_HOURS int `mapstructure:"JWT_EXPIRY_HOURS"`
	GIN_MODE         string `mapstructure:"GIN_MODE"`
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

	return cfg, nil
}
