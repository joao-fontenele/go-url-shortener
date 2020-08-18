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
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	})
	return logger
}
