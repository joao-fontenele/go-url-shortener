package router

import (
	"github.com/joao-fontenele/go-url-shortener/pkg/api/handler"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	routing "github.com/qiangxue/fasthttp-routing"
)

// New configures routes and it's handlers, and return it
func New(linkService shortener.LinkService) *routing.Router {
	router := routing.New()

	internalHandler := &handler.InternalHandler{}
	router.Get("/internal/status", internalHandler.StatusHandler)

	linkHandler := &handler.ShortenerHandler{LinkService: linkService}
	router.Post("/new", linkHandler.NewLink)
	router.Get("/<slug>", linkHandler.Redirect)

	return router
}
