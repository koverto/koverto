// Package resolver defines resolvers for GraphQL requests defined by the
// schema.
package resolver

import (
	"context"
	"fmt"

	koverto "github.com/koverto/koverto/api"

	"github.com/99designs/gqlgen/graphql"
	authz "github.com/koverto/authorization/api"
	"github.com/koverto/authorization/pkg/claims"
	credentials "github.com/koverto/credentials/api"
	"github.com/koverto/micro/v2"
	users "github.com/koverto/users/api"
	"github.com/koverto/uuid"
)

// Resolver defines a new set of GraphQL resolvers.
type Resolver struct {
	*micro.Service
	*micro.ClientSet
}

// New initializes a koverto.Config GraphQL service containing a Resolver.
func New() (*koverto.Config, error) {
	service, err := micro.NewService("com.koverto.svc.koverto", nil)
	if err != nil {
		return nil, err
	}

	clients := &micro.ClientSet{}
	clients.AddClient(authz.NewClient(service.Client()))
	clients.AddClient(credentials.NewClient(service.Client()))
	clients.AddClient(users.NewClient(service.Client()))

	return &koverto.Config{
		Resolvers: &Resolver{
			service,
			clients,
		},
		Directives: koverto.DirectiveRoot{
			Protected: protectedFieldDirective,
		},
	}, nil
}

// AuthorizationService returns the authorization service client belonging to
// the Resolver.
func (r *Resolver) AuthorizationService() authz.AuthorizationService {
	return r.ClientSet.Get(authz.Name).(authz.AuthorizationService)
}

// CredentialsService returns the credentials service client belonging to the
// Resolver.
func (r *Resolver) CredentialsService() credentials.CredentialsService {
	return r.ClientSet.Get(credentials.Name).(credentials.CredentialsService)
}

// UsersService returns the users service client belonging to the Resolver.
func (r *Resolver) UsersService() users.UsersService {
	return r.ClientSet.Get(users.Name).(users.UsersService)
}

func protectedFieldDirective(
	ctx context.Context,
	_ interface{},
	next graphql.Resolver,
	authRequired bool,
) (interface{}, error) {
	_, ok := ctx.Value(claims.ContextKeyJTI{}).(*uuid.UUID)

	if ok && !authRequired {
		return nil, fmt.Errorf("cannot do that while logged in")
	}

	if !ok && authRequired {
		return nil, fmt.Errorf("unauthorized")
	}

	return next(ctx)
}
