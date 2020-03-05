package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/koverto/micro"
	"github.com/rs/zerolog/log"
)

type logHandler struct {
	http.Handler
}

func LogHandler(next http.Handler) http.Handler {
	return &logHandler{next}
}

func (h *logHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()
	event := log.Info().Time("start_time", now)

	event.Str("http.method", r.Method)
	event.Str("http.path", r.URL.Path)
	event.Str("http.protocol", r.Proto)

	if referer := r.Referer(); referer != "" {
		event.Str("http.referer", referer)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if ip := net.ParseIP(host); err == nil && ip != nil {
		event.IPAddr("http.remote_addr", ip)
	}

	if rid, ok := micro.RequestIDFromContext(r.Context()); ok {
		event.Str("request_id", rid.Uuid.String())
	}

	lrw := &logResponseWriter{ResponseWriter: rw}
	h.Handler.ServeHTTP(lrw, r)

	event.Int("http.size", lrw.Size)
	event.Int("http.status", lrw.StatusCode)
	event.Dur("duration_ms", time.Since(now)).Send()
}

type logResponseWriter struct {
	http.ResponseWriter
	Size       int
	StatusCode int
}

func (rw *logResponseWriter) Write(b []byte) (int, error) {
	if rw.StatusCode == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.Size += size
	return size, err
}

func (rw *logResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}
