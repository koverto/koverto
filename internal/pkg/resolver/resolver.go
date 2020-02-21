// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	authz "github.com/koverto/authorization/api"
	credentials "github.com/koverto/credentials/api"
	koverto "github.com/koverto/koverto/api"
	users "github.com/koverto/users/api"
	"github.com/micro/go-micro/v2"
)

type Resolver struct {
	authz       authz.AuthorizationService
	credentials credentials.CredentialsService
	users       users.UsersService
}

func New() koverto.Config {
	service := micro.NewService(micro.Name("koverto"))
	service.Init()

	authz := authz.NewAuthorizationService("authorization", service.Client())
	credentials := credentials.NewCredentialsService("credentials", service.Client())
	users := users.NewUsersService("users", service.Client())

	return koverto.Config{
		Resolvers: &Resolver{
			authz,
			credentials,
			users,
		},
		Directives: koverto.DirectiveRoot{
			Protected: protectedFieldDirective,
		},
	}
}

func protectedFieldDirective(ctx context.Context, _ interface{}, next graphql.Resolver, authRequired bool) (interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}