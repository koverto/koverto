package middleware

import (
	"context"
	"net/http"
	"regexp"

	authorization "github.com/koverto/authorization/api"
	"github.com/koverto/authorization/pkg/claims"
	"github.com/koverto/koverto/internal/pkg/resolver"
)

const AUTHORIZATION_HEADER = "Authorization"

type authorizationHandler struct {
	*resolver.Resolver
	http.Handler
	bearerExpression *regexp.Regexp
}

func AuthorizationHandler(r *resolver.Resolver) func(http.Handler) http.Handler {
	bearerExpression, _ := regexp.Compile(`Bearer (\S+)`)
	return func(next http.Handler) http.Handler {
		return &authorizationHandler{
			Resolver:         r,
			Handler:          next,
			bearerExpression: bearerExpression,
		}
	}
}

func (h *authorizationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get(AUTHORIZATION_HEADER)

	if matches := h.bearerExpression.FindStringSubmatch(bearer); len(matches) == 2 {
		token := &authorization.Token{Token: matches[1]}

		if response, err := h.AuthorizationService.Validate(r.Context(), token); err == nil {
			ctx := context.WithValue(r.Context(), claims.ContextKeyUID{}, response.GetUserID().String())
			r = r.WithContext(ctx)
		}
	}

	h.Handler.ServeHTTP(rw, r)
}
