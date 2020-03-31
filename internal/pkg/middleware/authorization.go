// Package middleware defines middleware utilities for the user-facing GraphQL
// entrypoint.
package middleware

import (
	"context"
	"net/http"
	"regexp"

	"github.com/koverto/koverto/internal/pkg/resolver"

	authorization "github.com/koverto/authorization/api"
	"github.com/koverto/authorization/pkg/claims"
)

const authorizationHeader = "Authorization"
const bearerMatchLen = 2
const bearerPattern = `^Bearer (\S+)$`

type authorizationHandler struct {
	*resolver.Resolver
	http.Handler
	bearerExpression *regexp.Regexp
}

// AuthorizationHandler extracts and requests validation of a JWT present in the
// Authorization header of incoming HTTP requests.
func AuthorizationHandler(r *resolver.Resolver) func(http.Handler) http.Handler {
	bearerExpression := regexp.MustCompile(bearerPattern)

	return func(next http.Handler) http.Handler {
		return &authorizationHandler{
			Resolver:         r,
			Handler:          next,
			bearerExpression: bearerExpression,
		}
	}
}

func (h *authorizationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get(authorizationHeader)

	if matches := h.bearerExpression.FindStringSubmatch(bearer); len(matches) == bearerMatchLen {
		token := &authorization.Token{Token: matches[1]}

		if response, err := h.AuthorizationService().Validate(r.Context(), token); err == nil {
			ctx := context.WithValue(r.Context(), claims.ContextKeyJTI{}, response.GetID())
			ctx = context.WithValue(ctx, claims.ContextKeySUB{}, response.GetSubject())
			ctx = context.WithValue(ctx, claims.ContextKeyEXP{}, response.GetExpiresAt())
			r = r.WithContext(ctx)
		}
	}

	h.Handler.ServeHTTP(rw, r)
}
