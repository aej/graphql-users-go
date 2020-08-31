package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/andyjones11/graphql-users/graph/generated"
	"github.com/andyjones11/graphql-users/graph/model"
	"github.com/andyjones11/graphql-users/services/auth"
	userservice "github.com/andyjones11/graphql-users/services/user"
	"github.com/gin-gonic/gin"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*model.User, error) {
	http_response_ctx := ctx.Value("HttpResponse")
	http_response, _ := http_response_ctx.(*gin.Context)

	db_user, err := r.Repos.User.CreateUser(input.Email)
	auth.AuthenticateUser(r.Repos, db_user, http_response) // TODO: handle this error
	return DbUserToGqlUser(&db_user), err
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	httpResponse, _ := ctx.Value("HttpResponse").(*gin.Context)

	user, ok := ctx.Value("user").(userservice.User)

	if ok {
		auth.DeauthenticateUser(r.Repos, user, httpResponse)
		return true, nil
	}
	return false, nil
}

func (r *queryResolver) ListUsers(ctx context.Context) ([]*model.User, error) {
	db_users, _ := r.Repos.User.ListAllUsers()

	var users []*model.User

	for _, u := range db_users {
		users = append(users, DbUserToGqlUser(&u))
	}
	return users, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value("user").(userservice.User)

	if ok {
		return DbUserToGqlUser(&user), nil
	}
	return &model.User{}, errors.New("Authentication Required")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
