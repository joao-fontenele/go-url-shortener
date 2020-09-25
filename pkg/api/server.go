package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"github.com/joao-fontenele/go-url-shortener/pkg/postgres"
	"github.com/joao-fontenele/go-url-shortener/pkg/redis"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type statusResponse struct {
	Running bool `json:"running"`
}

func statusHandler(ctx *routing.Context) error {
	logger := logger.Get()
	logger.Info("GET /status")
	data := statusResponse{Running: true}
	b, _ := json.Marshal(data)
	ctx.SetContentType("application/json")
	ctx.Write(b)

	return nil
}

// Init sets up the server and it's routes
func Init() {
	err := configger.Load()
	if err != nil {
		log.Fatal("Failed to load configs: ", err)
	}

	logger := logger.Get()
	defer logger.Sync()

	dbClose, err := postgres.Connect()
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer dbClose()

	redisClose, err := redis.Connect()
	defer redisClose()

	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	router := routing.New()
	router.Get("/status", statusHandler)

	port := fmt.Sprintf(":%s", configger.Get().Port)
	logger.Sugar().Infof("listening on port %s", port)
	logger.Fatal("error", zap.Error(fasthttp.ListenAndServe(port, router.HandleRequest)))
}
