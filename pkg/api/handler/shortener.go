package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/joao-fontenele/go-url-shortener/pkg/api/response"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	routing "github.com/qiangxue/fasthttp-routing"
)

// NewLinkReqBody represents a request body received by the NewLink request handler
type NewLinkReqBody struct {
	URL string `json:"url"`
}

// ShortenerHandler is a route handler for link service
type ShortenerHandler struct {
	LinkService shortener.LinkService
}

// NewLink is a handler for creating a new Link
func (h *ShortenerHandler) NewLink(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	var body NewLinkReqBody
	err := json.Unmarshal(ctx.PostBody(), &body)

	if err != nil {
		status := http.StatusBadRequest
		ctx.SetStatusCode(status)
		b, _ := json.Marshal(response.HTTPErr{Message: "Invalid json in request body", StatusCode: status})
		ctx.Write(b)
		return nil
	}

	l, err := h.LinkService.Create(ctx, body.URL)
	if err != nil {
		status := http.StatusInternalServerError
		ctx.SetStatusCode(status)
		errMessage := fmt.Sprintf("Error creating link: %s", err.Error())
		b, _ := json.Marshal(response.HTTPErr{Message: errMessage, StatusCode: status})
		ctx.Write(b)
		return nil
	}

	ctx.SetStatusCode(http.StatusCreated)
	b, _ := json.Marshal(l)
	ctx.Write(b)

	return nil
}

// Redirect is a handler for redirecting to a Link.URL, given a slug from path
func (h *ShortenerHandler) Redirect(ctx *routing.Context) error {
	ctx.SetContentType("application/json")
	slug := ctx.Param("slug")
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
		return nil
	}

	ctx.Redirect(URL, http.StatusMovedPermanently)
	return nil
}
