package chilogger

import (
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type chilogger struct {
	logZ *zap.Logger
	name string
}

// NewZapMiddleware returns a new Zap Middleware handler.
func NewZapMiddleware(name string, logger *zap.Logger) func(next http.Handler) http.Handler {
	return chilogger{
		logZ: logger,
		name: name,
	}.middleware
}

func (c chilogger) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		} else {
			requestID = r.Header.Get("X-Request-Id")
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		if c.logZ != nil {
			fields := []zapcore.Field{
				zap.Int("http.response.status_code", ww.Status()),
				zap.Duration("event.duration", latency),
				zap.String("client.address", r.RemoteAddr),
				zap.String("url.path", r.RequestURI),
				zap.String("http.request.method", r.Method),
			}
			if requestID != "" {
				fields = append(fields, zap.String("http.request.id", requestID))
				fields = append(fields, zap.String("trace.id", requestID))
			}
			c.logZ.Info("request completed", fields...)
		}

	}
	return http.HandlerFunc(fn)
}
