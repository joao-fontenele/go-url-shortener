package router

import (
	"github.com/fasthttp/router"
	"github.com/joao-fontenele/go-url-shortener/pkg/api/handler"
	"github.com/joao-fontenele/go-url-shortener/pkg/api/middleware"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// New configures routes and it's handlers, and return it
func New(linkService shortener.LinkService) *router.Router {
	router := router.New()

	internalHandler := &handler.InternalHandler{}
	router.GET(
		"/internal/status",
		middleware.Logger(
			middleware.Metrics(internalHandler.StatusHandler),
		),
	)
	router.GET(
		"/internal/metrics",
		middleware.Logger(
			fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()),
		),
	)

	linkHandler := &handler.ShortenerHandler{LinkService: linkService}
	router.POST(
		"/links",
		middleware.Logger(
			middleware.Metrics(linkHandler.NewLink),
		),
	)
	router.GET(
		"/{slug}",
		middleware.Logger(
			middleware.Metrics(linkHandler.Redirect),
		),
	)

	return router
}
