package extmiddleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ZeroLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(wrappedWriter, r)

			log.WithLevel(zerolog.InfoLevel).
				Str("id", middleware.GetReqID(r.Context())).
				Str("remote_ip", r.RemoteAddr).
				Str("host", r.Host).
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Str("user_agent", r.UserAgent()).
				Int("status", wrappedWriter.Status()).
				Str("latency", time.Since(start).String()).
				Str("bytes_in", r.Header.Get("Content-Length")).
				Str("content_encoding", r.Header.Get("Content-Encoding")).
				Int64("bytes_out", int64(wrappedWriter.BytesWritten())).
				Msg("Chi server log")
		}
		return http.HandlerFunc(fn)
	}
}
