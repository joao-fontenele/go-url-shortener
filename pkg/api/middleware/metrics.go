package middleware

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"

	"github.com/joao-fontenele/go-url-shortener/pkg/metrics"
)

// Metrics is a middleware that handles default route metrics
func Metrics(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			statusCode := strconv.Itoa(ctx.Response.StatusCode())
			route := string(ctx.Path())

			metrics.HTTPRequestsCounter.With(
				prometheus.Labels{"route": route, "status": statusCode},
			).Inc()
			metrics.HTTPRequestsDurationHistogram.With(
				prometheus.Labels{"route": route, "status": statusCode},
			).Observe(time.Since(ctx.ConnTime()).Seconds())
		}()
		next(ctx)
		return
	}
}
