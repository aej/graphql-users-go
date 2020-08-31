package auth

import (
	"fmt"
	"github.com/andyjones11/graphql-users/repo"
	userservice "github.com/andyjones11/graphql-users/services/user"
	"math/rand"
)

// Creates a user session.
func AuthenticateUser(repo *repo.Repositories, user *userservice.User) (*userservice.User, error) {
	token := make([]byte, 32)
	rand.Read(token)

	user_token, _ := repo.User.CreateUserToken(user.ID, token)
	fmt.Printf("Created user token %s", user_token)

	return user, nil
}
