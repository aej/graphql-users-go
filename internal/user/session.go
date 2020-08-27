package userservice

import (
	"database/sql"
	"errors"
	"github.com/andyjones11/graphql-users/internal/storage"
	"github.com/google/uuid"
)

var SessionContext = "session"

// Create user session token for a given user
func CreateUserSessionToken(userRepo storage.UserRepo, user storage.User, token []byte) (storage.UsersToken, error) {
	sentTo := sql.NullString{Valid: false, String: ""}
	userToken, err := userRepo.CreateUserToken(user.ID, token, SessionContext, sentTo)
	if err != nil {
		return storage.UsersToken{}, errors.New("unable to create user session token")
	}
	return userToken, nil
}

func DestroyAllUserSessions(userRepo storage.UserRepo, userId uuid.UUID) {
	userRepo.DeleteUserToken(userId, SessionContext)
}