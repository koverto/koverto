// Package resolver defines resolvers for GraphQL requests defined by the
// schema.
package resolver

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	authz "github.com/koverto/authorization/api"
	"github.com/koverto/authorization/pkg/claims"
	credentials "github.com/koverto/credentials/api"
	koverto "github.com/koverto/koverto/api"
	"github.com/koverto/micro"
	users "github.com/koverto/users/api"
	"github.com/koverto/uuid"
)

// Resolver defines a new set of GraphQL resolvers.
type Resolver struct {
	authz.AuthorizationService
	credentials.CredentialsService
	users.UsersService
}

// New initializes a koverto.Config GraphQL service containing a Resolver.
func New() (*koverto.Config, error) {
	service, err := micro.NewService("com.koverto.svc.koverto", nil)
	if err != nil {
		return nil, err
	}

	authz := authz.NewAuthorizationService("authorization", service.Client())
	credentials := credentials.NewCredentialsService("credentials", service.Client())
	users := users.NewUsersService("users", service.Client())

	return &koverto.Config{
		Resolvers: &Resolver{
			authz,
			credentials,
			users,
		},
		Directives: koverto.DirectiveRoot{
			Protected: protectedFieldDirective,
		},
	}, nil
}

func protectedFieldDirective(ctx context.Context, _ interface{}, next graphql.Resolver, authRequired bool) (interface{}, error) {
	_, ok := ctx.Value(claims.ContextKeyJTI{}).(*uuid.UUID)

	if ok && !authRequired {
		return nil, fmt.Errorf("cannot do that while logged in")
	}

	if !ok && authRequired {
		return nil, fmt.Errorf("unauthorized")
	}

	return next(ctx)
}
