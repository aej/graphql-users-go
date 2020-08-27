package auth

import (
	"errors"
	"fmt"
	"github.com/andyjones11/graphql-users/internal/config"
	"github.com/andyjones11/graphql-users/internal/storage"
	userservice "github.com/andyjones11/graphql-users/internal/user"
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

var (
	InvalidCredentials = errors.New("invalid credentials")
	UserUnconfirmed    = errors.New("user unconfirmed")
	seed               *securecookie.SecureCookie
	secureSessionCookie bool
)

type AnonymousUser struct{}

func init() {
	var configuration = config.GetConfig()
	seed = securecookie.New([]byte(configuration.HashKey), []byte(configuration.SecretKey))
	secureSessionCookie = configuration.SecureSessionCookie
}

var sessionCookieName = "_id"
var sessionMaxAge = 60 * 60 * 24 * 60 // 60 days session lifespan

// Creates a user session in the session store and sets a cookie in the
// http response
func AuthenticateUser(repo *storage.Repositories, user storage.User, httpResponse echo.Context) (storage.User, error) {
	token := make([]byte, 32)
	rand.Read(token)

	_, err := userservice.CreateUserSessionToken(repo.User, user, token)

	if err != nil {
		return user, err
	}
	value, errs := seed.Encode(sessionCookieName, token)

	if errs != nil {
		fmt.Printf("%s", errs)
	}
	configuration := config.GetConfig()

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		MaxAge:   sessionMaxAge,
		Domain:   configuration.Domain,
		Secure:   secureSessionCookie,
		HttpOnly: true,
	}
	httpResponse.SetCookie(cookie)

	return user, nil
}

func DeauthenticateUser(repo *storage.Repositories, user storage.User, httpResponse echo.Context) {
	configuration := config.GetConfig()
	repo.User.DeleteUserToken(user.ID, userservice.SessionContext)

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		MaxAge:   0,
		Domain:   configuration.Domain,
		Secure:   configuration.SecureSessionCookie,
		HttpOnly: true,
	}
	httpResponse.SetCookie(cookie)
}

// Return the user from a cookie token
func CheckAuth(repo *storage.Repositories, httpResponse echo.Context) (storage.User, error) {
	sessionCookie, err := httpResponse.Cookie(sessionCookieName)
	if err != nil {
		return storage.User{}, err
	}
	token := make([]byte, 32)

	if err := seed.Decode(sessionCookieName, sessionCookie.Value, &token); err != nil {
		return storage.User{}, err
	}

	validUntil := time.Now().AddDate(0, 0, -60)
	userToken, err := repo.User.GetUnexpiredUserTokenForContext(token, userservice.SessionContext, validUntil)

	if err != nil {
		return storage.User{}, err
	}

	return userToken.User, nil
}

// Check a given email/password combination. One of three things can happen:
// 1) User with email does not exist
// 		- Return InvalidCredentials
// 2) User with email exists but password is incorrect
// 		- Return InvalidCredentials
// 3) User email has not been confirmed
// 		- Return UserUnconfirmed
func ValidateUserCredentials(repo *storage.Repositories, email string, password string) (storage.User, error) {
	user, err := repo.User.GetUserByEmail(email)

	if err != nil {
		return storage.User{}, InvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return storage.User{}, InvalidCredentials
	}
	if user.ConfirmedAt.Valid == false {
		return storage.User{}, UserUnconfirmed
	}
	return user, err
}
