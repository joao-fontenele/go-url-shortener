package configger

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config is a global variable that holds all the applications config
var configs Config

type database struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Name    string `mapstructure:"name"`
	User    string `mapstructure:"user"`
	Pass    string `mapstructure:"pass"`
	SSLMode string `mapstructure:"sslMode"`
}

type cache struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	CachePrefix     string `mapstructure:"cachePrefix"`
	LinksTTLSeconds int    `mapstructure:"linksTTLSeconds"`
}

// Config holds all applications configs
type Config struct {
	Env      string
	Port     string   `mapstructure:"port"`
	Database database `mapstructure:"database"`
	Cache    cache    `mapstructure:"cache"`
}

// Load configs from ./config/ yml files depending on APP_ENV.
// Order of precedence (local.yml > ${APP_ENV}.yml > default.yml)
func Load() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	v := viper.New()
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

// Get returns previously loaded configs
func Get() Config {
	return configs
}
