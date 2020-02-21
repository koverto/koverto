// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"
	"sync"

	authz "github.com/koverto/authorization/api"
	koverto "github.com/koverto/koverto/api"
	users "github.com/koverto/users/api"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input koverto.Authentication) (*koverto.LoginResponse, error) {
	user, err := r.users.Create(ctx, input.User)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 2)
	tokenCh := make(chan authz.Token, 1)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		input.Credential.UserID = user.GetId()
		_, err := r.credentials.Create(ctx, input.Credential)
		errCh <- err
	}()

	go func() {
		defer wg.Done()
		token, err := r.authz.Create(ctx, &authz.TokenRequest{
			UserID: user.GetId(),
		})
		errCh <- err
		tokenCh <- *token
	}()

	wg.Wait()
	close(errCh)
	close(tokenCh)

	for err = range errCh {
		if err != nil {
			return nil, err
		}
	}

	token := <-tokenCh

	return &koverto.LoginResponse{
		Token: token.GetToken(),
		User:  user,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input koverto.Authentication) (*koverto.LoginResponse, error) {
	user, err := r.users.Read(ctx, input.User)
	if err != nil {
		return nil, err
	}

	input.Credential.UserID = user.GetId()
	if _, err := r.credentials.Validate(ctx, input.Credential); err != nil {
		return nil, err
	}

	token, err := r.authz.Create(ctx, &authz.TokenRequest{
		UserID: user.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &koverto.LoginResponse{
		Token: token.GetToken(),
		User:  user,
	}, nil
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
