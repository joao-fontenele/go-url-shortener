package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
)

var rdb *redis.Client

// Connect creates a connection to redis
func Connect() (func() error, error) {
	logger := logger.Get()
	dbConf := configger.Get().Cache
	logger.Info("connecting to redis")

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", dbConf.Host, dbConf.Port),
	})

	_, err := rdb.Ping(context.Background()).Result()
	return rdb.Close, err
}

// GetConnection returns a previously created connection pool
func GetConnection() *redis.Client {
	return rdb
}
