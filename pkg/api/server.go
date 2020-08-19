package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joao-fontenele/go-url-shortener/pkg/common"
	"go.uber.org/zap"
)

type statusResponse struct {
	Running bool `json:"running"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	data := statusResponse{Running: true}
	js, err := json.Marshal(data)
	logger := common.GetLogger()
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
	logger := common.GetLogger()
	defer logger.Sync()
	http.HandleFunc("/status", statusHandler)
	logger.Info("listening on port 8080")
	logger.Fatal("error", zap.Error(http.ListenAndServe(":8080", nil)))
}
