package logger

import (
	"sync"

	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"go.uber.org/zap"
)

var logger *zap.Logger
var once sync.Once

// Get gets singleton logger instance
func Get() *zap.Logger {
	once.Do(func() {
		var err error

		env := configger.Get().Env
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
