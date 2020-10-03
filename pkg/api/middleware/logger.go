package middleware

import (
	"time"

	"github.com/joao-fontenele/go-url-shortener/pkg/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Logger is a middleware that logs request details
func Logger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		l := logger.Get()

		reqID := zap.Uint64("reqId", ctx.ID())
		localAddr := zap.String("localAddr", ctx.LocalAddr().String())
		remoteAddr := zap.String("remoteAddr", ctx.RemoteAddr().String())
		method := zap.ByteString("method", ctx.Method())
		uri := zap.ByteString("uri", ctx.URI().RequestURI())

		l.Debug(
			"Request received",
			reqID,
			localAddr,
			remoteAddr,
			method,
			uri,
		)
		defer func() {
			elapsed := time.Since(ctx.ConnTime()).Seconds()
			l.Info(
				"Request ended",
				reqID,
				localAddr,
				remoteAddr,
				method,
				uri,
				zap.Float64("elapsedSeconds", elapsed),
			)
		}()

		next(ctx)
		return
	}
}
