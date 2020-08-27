package graph

import (
	"github.com/andyjones11/graphql-users/internal/graph/model"
	"github.com/andyjones11/graphql-users/internal/storage"
)

func DbUserToGqlUser(user *storage.User) *model.User {
	u := &model.User{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.CreatedAt.String(),
	}
	return u
}
