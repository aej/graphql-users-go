package repo

import (
	"errors"
	"github.com/andyjones11/graphql-users/services/user"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

var UserEmailExists = errors.New("user email exists")

func (r *UserRepo) GetUserById(id string) (userservice.User, error) {
	var user userservice.User

	result := r.db.Where("id = ?", id).First(&user)

	if result.RowsAffected == 1 {
		return user, nil
	}
	return user, RecordNotFound
}

func (r *UserRepo) CreateUser(email string) (userservice.User, error) {
	user := userservice.User{Email: email, HashedPassword: "something really secret"}

	var existingUser userservice.User

	existingResult := r.db.Where("email = ?", user.Email).First(&existingUser)

	if existingResult.RowsAffected == 1 {
		return user, UserEmailExists
	}

	r.db.Create(&user)
	return user, nil
}

func (r *UserRepo) ListAllUsers() ([]userservice.User, error) {
	var users []userservice.User
	r.db.Find(&users)
	return users, nil
}

func (r *UserRepo) CreateUserToken(UserId uuid.UUID, token []byte) (userservice.UsersToken, error) {
	userToken := userservice.UsersToken{UserId: UserId, Token: token}
	r.db.Create(&userToken)

	return userToken, nil
}

// Delete All UsersTokens for an associated User
func (r *UserRepo) DeleteUserToken(UserId uuid.UUID) {
	var userTokens []userservice.UsersToken
	result := r.db.Where("user_id = ?", UserId).Find(&userTokens)

	if result.RowsAffected > 0 {
		var idsToDelete []string
		for _, token := range userTokens {
			idsToDelete = append(idsToDelete, token.ID.String())
		}
		r.db.Delete(&userTokens, idsToDelete) // delete all tokens with this primary key
	}
}

var RecordNotFound = errors.New("Record could not be found")

func (r *UserRepo) GetUserTokenByToken(token []byte) (userservice.UsersToken, error) {
	var usersToken userservice.UsersToken
	result := r.db.Where("token = ?", token).First(&usersToken)
	if result.RowsAffected == 1 {
		return usersToken, nil
	}
	return usersToken, RecordNotFound
}

type Repositories struct {
	User UserRepo
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	repositories := &Repositories{
		User: UserRepo{db: db},
	}
	return repositories, nil
}
