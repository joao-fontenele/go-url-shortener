package main

import (
	"fmt"

	"github.com/joao-fontenele/go-url-shortener/pkg/api"
	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	router := api.New()

	logger := logger.Get()

	port := fmt.Sprintf(":%s", configger.Get().Port)
	logger.Sugar().Infof("listening on port %s", port)
	logger.Fatal("error", zap.Error(fasthttp.ListenAndServe(port, router.HandleRequest)))
}
