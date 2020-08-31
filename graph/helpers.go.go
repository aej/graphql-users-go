package graph

import (
	"github.com/andyjones11/graphql-users/graph/model"
	"github.com/andyjones11/graphql-users/services/user"
)

func DbUserToGqlUser(user *userservice.User) *model.User {
	u := &model.User{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.CreatedAt.String(),
	}
	return u
}
