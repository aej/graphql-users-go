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

func (r *UserRepo) CreateUser(email string) (userservice.User, error) {
	user := userservice.User{Email: email, HashedPassword: "something really secret"}

	var existing_user userservice.User

	existing_result := r.db.Where("email = ?", user.Email).First(&existing_user)

	if existing_result.RowsAffected == 1 {
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
	user_token := userservice.UsersToken{UserId: UserId, Token: token}
	r.db.Create(&user_token)

	return user_token, nil
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
