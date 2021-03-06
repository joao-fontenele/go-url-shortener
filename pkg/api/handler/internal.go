package handler

import (
	"encoding/json"

	"github.com/joao-fontenele/go-url-shortener/pkg/api/response"
	"github.com/valyala/fasthttp"
)

// InternalHandler will handle internal, non public requests
type InternalHandler struct {
}

// TODO: handle concurrent access?
var staticResponseCache []byte

// StatusHandler handles requests to api's status route
func (h *InternalHandler) StatusHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	if len(staticResponseCache) == 0 {
		data := response.StatusResponse{Running: true}
		staticResponseCache, _ = json.Marshal(data)
	}

	ctx.Write(staticResponseCache)
}
