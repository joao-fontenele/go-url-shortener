package common

import (
	"sync"

	"go.uber.org/zap"
)

var logger *zap.Logger
var once sync.Once

// GetLogger gets singleton logger instance
func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error

		env := GetConf().Env
		if env == "development" {
			logger, err = zap.NewDevelopment()
		} else if env == "test" {
			logger = zap.NewNop()
		} else {
			logger, err = zap.NewProduction()
		}

		if err != nil {
			panic(err)
		}
	})

	return logger
}
