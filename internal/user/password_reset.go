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

var passwordResetContext = "password_reset"

// Request password reset for a user email. If the User exists
// then triggers a password reset email
func RequestResetPassword(userRepo storage.UserRepo, email string) error {
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user with email does not exist")
	}

	if err := TriggerPasswordResetEmail(userRepo, user); err != nil {
		return err
	}

	return nil
}

// Checks the validity of a password reset token. If valid then return the
// associated User
func ValidatePasswordResetToken(userRepo storage.UserRepo, token string) (storage.User, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return storage.User{}, errors.New("invalid token")
	}

	hashedToken := sha256.New()
	hashedToken.Write(decodedToken)

	validUntil := time.Now().AddDate(0, 0, -1)
	userToken, err := userRepo.GetUnexpiredUserTokenForContext(hashedToken.Sum(nil), passwordResetContext, validUntil)

	if err != nil {
		return storage.User{}, errors.New("invalid token")
	}

	return userToken.User, nil
}

// Reset a user password. Validate that the password reset token
// exists and has not expired.
// Once password has successfully been reset all active sessions are destroyed for the User
// forcing the user to authenticate again.
func ResetPassword(userRepo storage.UserRepo, password string, token string) (bool, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return false, errors.New("invalid token")
	}

	hashedToken := sha256.New()
	hashedToken.Write(decodedToken)

	validUntil := time.Now().AddDate(0, 0, -1)
	usersToken, err := userRepo.GetUnexpiredUserTokenForContext(hashedToken.Sum(nil), passwordResetContext, validUntil)

	if err != nil {
		return false, errors.New("invalid token")
	}

	user, err := userRepo.UpdateUserPassword(usersToken.User, password)

	if err != nil {
		return false, errors.New("unable to reset password")
	}

	userRepo.DeleteUserToken(user.ID, passwordResetContext)
	DestroyAllUserSessions(userRepo, user.ID)

	return true, nil
}

// Create a password reset token and send the user a password reset email
func TriggerPasswordResetEmail(userRepo storage.UserRepo, user storage.User) error {
	// create 32 random bytes
	token := make([]byte, 32)
	rand.Read(token)

	userToken, err := createPasswordResetToken(userRepo, user, token)

	if err != nil {
		panic("unable to create password reset token for user")
	}

	passwordResetToken := base64.StdEncoding.EncodeToString(token)
	email.SendPasswordResetEmail(userToken.SentTo.String, passwordResetToken)

	return nil
}

// Create password reset token for a given user
func createPasswordResetToken(userRepo storage.UserRepo, user storage.User, token []byte) (storage.UsersToken, error) {
	hashedToken := sha256.New()
	hashedToken.Write(token)
	sentTo := sql.NullString{Valid: true, String: user.Email}
	userToken, err := userRepo.CreateUserToken(user.ID, hashedToken.Sum(nil), passwordResetContext, sentTo)

	if err != nil {
		return userToken, errors.New("unable to create email confirmation token")
	}
	return userToken, nil
}
