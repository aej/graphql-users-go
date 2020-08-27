package graph

import (
	"context"
	"github.com/andyjones11/graphql-users/graph/generated"
	"github.com/andyjones11/graphql-users/graph/model"
	"github.com/andyjones11/graphql-users/repo"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterInput) (*model.User, error) {
	user := repo.User{Email: input.Email}
	return &model.User{
		ID: user.ID.String(),
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	db_users, _ := r.Repo.User.ListAllUsers()

	var users []*model.User

	for _, u := range db_users {
		users = append(users, &model.User{
			ID: u.ID.String(),
			Email: u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
