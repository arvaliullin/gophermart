package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

const (
	msgRequestProcessed = "запрос обработан"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Logging создаёт middleware для логирования HTTP запросов.
func Logging(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logger.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Int("status", wrapped.statusCode).
				Int("size", wrapped.size).
				Dur("duration", duration).
				Msg(msgRequestProcessed)
		})
	}
}

