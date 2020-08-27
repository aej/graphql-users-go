package userservice

import (
	"github.com/andyjones11/graphql-users/internal/storage"
)

// Create a new user. Send email confirmation to user if successful
func CreateUser(userRepo storage.UserRepo, email string, password string, fullName string) (storage.User, error) {
	dbUser, err := userRepo.CreateUser(email, password, fullName)

	if err := TriggerEmailConfirmation(userRepo, dbUser); err != nil {
		return dbUser, err
	}

	return dbUser, err
}

