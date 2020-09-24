package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"github.com/joao-fontenele/go-url-shortener/pkg/postgres"
	"github.com/joao-fontenele/go-url-shortener/pkg/redis"
	"go.uber.org/zap"
)

type statusResponse struct {
	Running bool `json:"running"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	data := statusResponse{Running: true}
	js, err := json.Marshal(data)
	logger := logger.Get()
	logger.Info("GET /status")
	if err != nil {
		fmt.Println("err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

	http.HandleFunc("/status", statusHandler)

	port := fmt.Sprintf(":%s", configger.Get().Port)

	logger.Sugar().Infof("listening on port %s", port)
	logger.Fatal("error", zap.Error(http.ListenAndServe(port, nil)))
}
