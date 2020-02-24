// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	authz "github.com/koverto/authorization/api"
	credentials "github.com/koverto/credentials/api"
	koverto "github.com/koverto/koverto/api"
	"github.com/koverto/micro"
	users "github.com/koverto/users/api"
)

type Resolver struct {
	authz.AuthorizationService
	credentials.CredentialsService
	users.UsersService
}

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
	return next(ctx)
	// panic(fmt.Errorf("not implemented"))
}
