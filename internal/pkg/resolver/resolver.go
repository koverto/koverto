// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	koverto "github.com/koverto/koverto/api"
)

type Resolver struct{}

func New() koverto.Config {
	return koverto.Config{
		Resolvers: &Resolver{},
		Directives: koverto.DirectiveRoot{
			Protected: protectedFieldDirective,
		},
	}
}

func protectedFieldDirective(ctx context.Context, _ interface{}, next graphql.Resolver, authRequired bool) (interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}
