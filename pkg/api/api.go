package api

import (
	"log"

	"github.com/joao-fontenele/go-url-shortener/pkg/api/router"
	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"github.com/joao-fontenele/go-url-shortener/pkg/postgres"
	"github.com/joao-fontenele/go-url-shortener/pkg/redis"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	routing "github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"
)

func loadConfs() {
	err := configger.Load()
	if err != nil {
		log.Fatalf("Failed to load configs: %v", err)
	}
}

func connectDB(logger *zap.Logger) {
	_, err := postgres.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
}

func connectCache(logger *zap.Logger) {
	_, err := redis.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to redis", zap.Error(err))
	}
}

func newLinkService() shortener.LinkService {
	dbConn := postgres.GetConnection()
	dbDao := postgres.NewLinkDao(dbConn)

	cacheConn := redis.GetConnection()
	cacheDao := redis.NewLinkDao(cacheConn)

	linkRepo := shortener.NewLinkRepository(dbDao, cacheDao)

	return shortener.NewLinkService(linkRepo)
}

// New loads configs, sets up connection, and api routes
func New() *routing.Router {
	loadConfs()

	logger := logger.Get()

	connectDB(logger)
	connectCache(logger)

	ls := newLinkService()
	router := router.New(ls)

	return router
}