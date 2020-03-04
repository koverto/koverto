package middleware

import (
	"net/http"

	"github.com/koverto/micro"
	"github.com/koverto/uuid"
)

const REQUEST_ID_HEADER = "X-Request-Id"

type ContextKeyRequestID struct{}

type requestIDHandler struct {
	http.Handler
}

func RequestIDHandler(next http.Handler) http.Handler {
	return &requestIDHandler{next}
}

func (h *requestIDHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rid := uuid.New()
	rw.Header().Add(REQUEST_ID_HEADER, rid.Uuid.String())

	ctx := micro.ContextWithRequestID(r.Context(), rid)
	h.Handler.ServeHTTP(rw, r.WithContext(ctx))
}
