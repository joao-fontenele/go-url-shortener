package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/joao-fontenele/go-url-shortener/pkg/api/response"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	"github.com/valyala/fasthttp"
)

// newLinkReqBody represents a request body received by the NewLink request handler
type newLinkReqBody struct {
	URL string `json:"url"`
}

// ShortenerHandler is a route handler for link service
type ShortenerHandler struct {
	LinkService shortener.LinkService
}

// NewLink is a handler for creating a new Link
func (h *ShortenerHandler) NewLink(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	var body newLinkReqBody
	err := json.Unmarshal(ctx.PostBody(), &body)

	if err != nil {
		status := http.StatusBadRequest
		ctx.SetStatusCode(status)
		b, _ := json.Marshal(response.HTTPErr{Message: "Invalid json in request body", StatusCode: status})
		ctx.Write(b)
		return
	}

	l, err := h.LinkService.Create(ctx, body.URL)
	if err != nil {
		var status int
		var errMessage string

		if errors.Is(err, shortener.ErrInvalidLink) {
			status = http.StatusBadRequest
			errMessage = err.Error()
		} else {
			status = http.StatusInternalServerError
			errMessage = fmt.Sprintf("Error creating link: %s", err.Error())
		}

		b, _ := json.Marshal(response.HTTPErr{Message: errMessage, StatusCode: status})
		ctx.SetStatusCode(status)
		ctx.Write(b)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
	b, _ := json.Marshal(l)
	ctx.Write(b)

	return
}

// Redirect is a handler for redirecting to a Link.URL, given a slug from path
func (h *ShortenerHandler) Redirect(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))
	URL, err := h.LinkService.GetURL(ctx, slug)
	if err != nil {
		var status int
		var errMessage string

		if errors.Is(err, shortener.ErrLinkNotFound) {
			status = http.StatusNotFound
			errMessage = fmt.Sprintf("Link with slug '%s' not found", slug)
		} else {
			status = http.StatusInternalServerError
			errMessage = fmt.Sprintf("Error getting slug: %s", err.Error())
		}

		ctx.SetStatusCode(status)
		b, _ := json.Marshal(response.HTTPErr{Message: errMessage, StatusCode: status})
		ctx.Write(b)
		return
	}

	ctx.Redirect(URL, http.StatusMovedPermanently)
	return
}
