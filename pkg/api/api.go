package api

import (
	"log"

	"github.com/fasthttp/router"
	myRouter "github.com/joao-fontenele/go-url-shortener/pkg/api/router"
	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"github.com/joao-fontenele/go-url-shortener/pkg/metrics"
	"github.com/joao-fontenele/go-url-shortener/pkg/postgres"
	"github.com/joao-fontenele/go-url-shortener/pkg/redis"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
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

func initMetrics() {
	metrics.Init()
}

// New loads configs, sets up connection, and api routes
func New() *router.Router {
	loadConfs()

	logger := logger.Get()

	connectDB(logger)
	connectCache(logger)

	initMetrics()

	ls := newLinkService()
	r := myRouter.New(ls)

	return r
}
