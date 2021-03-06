package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sync"
	"time"

	authz "github.com/koverto/authorization/api"
	"github.com/koverto/authorization/pkg/claims"
	koverto "github.com/koverto/koverto/api"
	users "github.com/koverto/users/api"
	"github.com/koverto/uuid"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input koverto.Authentication) (*koverto.LoginResponse, error) {
	user, err := r.UsersService().Create(ctx, input.User)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error)
	tokenCh := make(chan authz.Token, 1)

	wg := sync.WaitGroup{}
	wgFns := make([]func(), 0)

	wgFns = append(wgFns, func() {
		defer wg.Done()

		input.Credential.UserID = user.GetId()
		_, err := r.CredentialsService().Create(ctx, input.Credential)
		errCh <- err
	})

	wgFns = append(wgFns, func() {
		defer wg.Done()

		token, err := r.AuthorizationService().Create(ctx, &authz.Claims{
			Subject: user.GetId(),
		})
		errCh <- err
		tokenCh <- *token
	})

	wg.Add(len(wgFns))

	for _, fn := range wgFns {
		go fn()
	}

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
	loginFailed := fmt.Errorf("invalid e-mail address or password")

	user, err := r.UsersService().Read(ctx, input.User)
	if err != nil {
		return nil, loginFailed
	}

	input.Credential.UserID = user.GetId()
	if _, err := r.CredentialsService().Validate(ctx, input.Credential); err != nil {
		return nil, loginFailed
	}

	token, err := r.AuthorizationService().Create(ctx, &authz.Claims{
		Subject: user.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &koverto.LoginResponse{
		Token: token.GetToken(),
		User:  user,
	}, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (*koverto.LogoutResponse, error) {
	jti, ok := ctx.Value(claims.ContextKeyJTI{}).(*uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("no token ID")
	}

	exp, ok := ctx.Value(claims.ContextKeyEXP{}).(*time.Time)
	if !ok {
		return nil, fmt.Errorf("no token expiry")
	}

	claims := &authz.Claims{ID: jti, ExpiresAt: exp}
	_, err := r.AuthorizationService().Invalidate(ctx, claims)

	return &koverto.LogoutResponse{Ok: err == nil}, err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input users.User) (*users.User, error) {
	uid, ok := ctx.Value(claims.ContextKeySUB{}).(*uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("no user ID")
	}

	input.Id = uid

	return r.UsersService().Update(ctx, &input)
}

func (r *queryResolver) GetUser(ctx context.Context) (*users.User, error) {
	uid, ok := ctx.Value(claims.ContextKeySUB{}).(*uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("no user ID")
	}

	return r.UsersService().Read(ctx, &users.User{
		Id: uid,
	})
}

// Mutation returns koverto.MutationResolver implementation.
func (r *Resolver) Mutation() koverto.MutationResolver { return &mutationResolver{r} }

// Query returns koverto.QueryResolver implementation.
func (r *Resolver) Query() koverto.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
