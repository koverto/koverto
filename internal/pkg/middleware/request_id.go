package middleware

import (
	"net/http"

	"github.com/koverto/micro"
	"github.com/koverto/uuid"
)

const requestIDHeader = "X-Request-Id"

type requestIDHandler struct {
	http.Handler
}

// RequestIDHandler appends a request ID to the context of incoming requests and
// the headers of the accompanying response.
func RequestIDHandler(next http.Handler) http.Handler {
	return &requestIDHandler{next}
}

func (h *requestIDHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rid := uuid.New()
	rw.Header().Add(requestIDHeader, rid.Uuid.String())

	ctx := micro.ContextWithRequestID(r.Context(), rid)
	h.Handler.ServeHTTP(rw, r.WithContext(ctx))
}
