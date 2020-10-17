package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"go.uber.org/zap"
)

var rdb *redis.Client

// Connect creates a connection to redis
func Connect() (func() error, error) {
	logger := logger.Get()
	dbConf := configger.Get().Cache
	logger.Info("Connecting to redis")

	var connectURL string
	if dbConf.ConnectURL != "" {
		connectURL = dbConf.ConnectURL
	} else {
		connectURL = fmt.Sprintf("%s:%s", dbConf.Host, dbConf.Port)
	}

	options, err := redis.ParseURL(connectURL)
	options.Username = "" // workaround for connecting to a redis server 5
	if err != nil {
		logger.Fatal("Failed to parse redis connection string", zap.Error(err))
	}

	rdb = redis.NewClient(options)

	_, err = rdb.Ping(context.Background()).Result()
	return rdb.Close, err
}

// GetConnection returns a previously created connection pool
func GetConnection() *redis.Client {
	return rdb
}
