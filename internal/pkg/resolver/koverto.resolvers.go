// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	koverto "github.com/koverto/koverto/api"
	users "github.com/koverto/users/api"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input koverto.Authentication) (*koverto.LoginResponse, error) {
	user, err := r.users.Create(ctx, input.User)
	if err != nil {
		return nil, err
	}

	input.Credential.UserID = user.GetId()
	if _, err := r.credentials.Create(ctx, input.Credential); err != nil {
		return nil, err
	}

	return &koverto.LoginResponse{
		Token: "token goes here",
		User:  user,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input koverto.Authentication) (*koverto.LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input users.User) (*users.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUser(ctx context.Context) (*users.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Mutation() koverto.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() koverto.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
