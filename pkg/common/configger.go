package common

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config is a global variable taht holds all the applications config
var configs Config

type database struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Name    string `mapstructure:"name"`
	User    string `mapstructure:"user"`
	Pass    string `mapstructure:"pass"`
	SSLMode string `mapstructure:"sslMode"`
}

// Config holds all apllications configs
type Config struct {
	Env      string
	Port     string   `mapstructure:"port"`
	Database database `mapstructure:"database"`
}

// LoadConfs configs from ./config/ yml files depending on APP_ENV.
// Order of precedence (local.yml > ${APP_ENV}.yml > default.yml)
func LoadConfs() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	v := viper.New()
	v.SetEnvPrefix("APP")
	v.Set("env", env)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	v.SetConfigType("yaml")

	// Load default configs file
	v.SetConfigName("config/default")
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// Override with env specific configs
	v.SetConfigName(fmt.Sprintf("config/%s", env))
	if err := v.MergeInConfig(); err != nil {
		return err
	}

	// Override with local config file. This file is not tracked and a safe to place secrets
	v.SetConfigName("config/local")
	v.MergeInConfig() // ignore missing local file

	v.Unmarshal(&configs)
	return nil
}

// GetConf returns previously loaded configs
func GetConf() Config {
	return configs
}
