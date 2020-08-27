package graph

import (
	"errors"
	"github.com/vektah/gqlparser/gqlerror"
)

var UserEmailExists = errors.New("User with email already exists")

var AuthenticationRequired = gqlerror.Errorf("Authentication Required")
var AuthenticatedDisallowed = gqlerror.Errorf("Authenticated users are disallowed")
var InvalidEmailPasswordCombination = gqlerror.Errorf("Invalid Email/Password combination")
var UnconfirmedEmail = gqlerror.Errorf("Email is unconfirmed")
var ServerError = errors.New("A server error occurred")
