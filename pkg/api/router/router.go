package router

import (
	"github.com/fasthttp/router"
	"github.com/joao-fontenele/go-url-shortener/pkg/api/handler"
	"github.com/joao-fontenele/go-url-shortener/pkg/api/middleware"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

// New configures routes and it's handlers, and return it
func New(linkService shortener.LinkService) *router.Router {
	router := router.New()

	internalHandler := &handler.InternalHandler{}
	router.GET("/internal/status", middleware.Logger(internalHandler.StatusHandler))

	linkHandler := &handler.ShortenerHandler{LinkService: linkService}
	router.POST("/links", middleware.Logger(linkHandler.NewLink))
	router.GET("/{slug}", middleware.Logger(linkHandler.Redirect))

	return router
}
