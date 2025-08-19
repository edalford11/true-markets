package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Environment string

const (
	Production  Environment = "production"
	Staging     Environment = "staging"
	Development Environment = "development"
	Test        Environment = "test"
	Local       Environment = "local"
)

func IsTestEnvironment() bool {
	return strings.ToLower(viper.GetString("environment")) == string("test")
}
func IsLocalEnvironment() bool {
	return strings.ToLower(viper.GetString("environment")) == string(Local)
}

// Init initializes the viper config singleton.
func Init() {
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	env := os.Getenv("ENVIRONMENT")

	if env == "test" {
		viper.SetConfigName("test.default")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("test")
		err = viper.MergeInConfig()
		if err != nil {
			if errors.As(err, &viper.ConfigFileNotFoundError{}) {
				slog.Debug("test.yml not found using test.default.yml")
				return
			}
			panic(err)
		}
	} else {
		viper.SetConfigName("config.default")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("config")
		err = viper.MergeInConfig()
		if err != nil {
			if errors.As(err, &viper.ConfigFileNotFoundError{}) {
				slog.Debug("config.yml not found using config.default.yml")
				return
			}
			panic(err)
		}
	}
}
