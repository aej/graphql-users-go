package auth

import (
	"fmt"
	"github.com/andyjones11/graphql-users/repo"
	userservice "github.com/andyjones11/graphql-users/services/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"math/rand"
)

var (
	hashKey   = []byte("andyjones1111111")
	secretKey = []byte("andyjones1111111") // should be set by the environment and must be 16, 24 or 32 byte strings
	seed      = securecookie.New(hashKey, secretKey)
)

var sessionCookieName = "_id"

// Creates a user session in the session store and sets a cookie in the
// http response
func AuthenticateUser(repo *repo.Repositories, user userservice.User, httpResponse *gin.Context) (userservice.User, error) {
	token := make([]byte, 32)
	rand.Read(token)

	// TODO: handle this error
	repo.User.CreateUserToken(user.ID, token)

	v, errs := seed.Encode(sessionCookieName, token)

	if errs != nil {
		fmt.Printf("%s", errs)
	}
	httpResponse.SetCookie(sessionCookieName, v, 3600, "/", "localhost", false, true)

	return user, nil
}

func DeauthenticateUser(repo *repo.Repositories, user userservice.User, httpResponse *gin.Context) {
	repo.User.DeleteUserToken(user.ID)
	httpResponse.SetCookie(sessionCookieName, "", 0, "/", "localhost", false, true)
}

// Return the user from a cookie token
func CheckAuth(repo *repo.Repositories, httpResponse *gin.Context) (userservice.User, error) {
	sessionCookie, err := httpResponse.Cookie(sessionCookieName)
	if err != nil {
		return userservice.User{}, err
	}
	token := make([]byte, 32)
	if err := seed.Decode(sessionCookieName, sessionCookie, &token); err != nil {
		return userservice.User{}, err
	}
	userToken, err := repo.User.GetUserTokenByToken(token)
	if err != nil {
		return userservice.User{}, err
	}
	user, err := repo.User.GetUserById(userToken.UserId.String())
	if err != nil {
		return userservice.User{}, err
	}
	return user, nil
}
