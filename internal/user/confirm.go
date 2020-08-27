package userservice

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/andyjones11/graphql-users/internal/email"
	"github.com/andyjones11/graphql-users/internal/storage"
	"math/rand"
	"time"
)

var emailConfirmationContext = "confirm"

var UserEmailAlreadyConfirmed = errors.New("user email already confirmed")

// Create email confirmation token for a given user
func CreateEmailConfirmationToken(userRepo storage.UserRepo, user storage.User, token []byte) (storage.UsersToken, error) {
	hashedToken := sha256.New()
	hashedToken.Write(token)
	sendTo := sql.NullString{Valid: true, String: user.Email}
	userToken, err := userRepo.CreateUserToken(user.ID, hashedToken.Sum(nil), emailConfirmationContext, sendTo)

	if err != nil {
		return userToken, errors.New("unable to create email confirmation token")
	}
	return userToken, nil
}

var TokenNotFound = errors.New("token not found")
var InvalidToken = errors.New("invalid confirmation token")

// Given a confirmation token attempt to confirm the user email
// Returns the following errors:
// 1) TokenNotFound
// 2) InvalidToken
func ConfirmUserEmail(userRepo storage.UserRepo, confirmationToken string) error {
	decodedToken, err := base64.StdEncoding.DecodeString(confirmationToken)

	if err != nil {
		return InvalidToken
	}

	hashedToken := sha256.New()
	hashedToken.Write(decodedToken)

	validUntil := time.Now().AddDate(0, 0, -7)
	usersToken, err := userRepo.GetUnexpiredUserTokenForContext(hashedToken.Sum(nil), emailConfirmationContext, validUntil)
	if err != nil {
		return TokenNotFound
	}
	userRepo.ConfirmUserEmail(usersToken.User)
	userRepo.DeleteUserToken(usersToken.UserId, emailConfirmationContext)

	return nil
}

func TriggerEmailConfirmation(userRepo storage.UserRepo, user storage.User) error {
	if user.ConfirmedAt.Valid == true {
		return UserEmailAlreadyConfirmed
	}

	// create 32 random bytes
	token := make([]byte, 32)
	rand.Read(token)

	_, err := CreateEmailConfirmationToken(userRepo, user, token)

	if err != nil {
		return err
	}

	confirmationToken := base64.StdEncoding.EncodeToString(token)
	email.SendConfirmationEmail(user.Email, confirmationToken)

	return nil
}
