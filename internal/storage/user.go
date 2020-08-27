package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             uuid.UUID    `gorm:"primaryKey;type:uuid"`
	Email          string       `gorm:"type:varchar(255);unique_index;not null"`
	FullName       string       `gorm:"type:varchar(255);not null"`
	HashedPassword string       `gorm:"not null"`
	ConfirmedAt    sql.NullTime `gorm:"null"`
	CreatedAt      time.Time    `gorm:"not null"`
	UpdatedAt      time.Time    `gorm:"not null"`
	UserTokens     []UsersToken
}

type UserRepo struct {
	db *gorm.DB
}

var UserEmailExists = errors.New("user email exists")

func (r *UserRepo) GetUserById(id string) (User, error) {
	var user User

	result := r.db.Where("id = ?", id).First(&user)

	if result.RowsAffected == 1 {
		return user, nil
	}
	return user, RecordNotFound
}

func (r *UserRepo) GetUserByEmail(email string) (User, error) {
	var user User

	result := r.db.Where("email = ?", email).First(&user)

	if result.RowsAffected == 1 {
		return user, nil
	}
	return user, RecordNotFound
}

func (r *UserRepo) CreateUser(email string, password string, fullName string) (User, error) {
	user := User{
		Email:    email,
		FullName: fullName,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic("unable to hash password")
	}

	user.HashedPassword = string(hashedPassword)

	var existingUser User

	existingResult := r.db.Where("email = ?", user.Email).First(&existingUser)

	if existingResult.RowsAffected == 1 {
		return user, UserEmailExists
	}

	r.db.Create(&user)
	return user, nil
}

func (r *UserRepo) UpdateUserPassword(user User, password string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic("unable to hash password")
	}
	user.HashedPassword = string(hashedPassword)
	r.db.Save(&user)

	return user, nil
}

func (r *UserRepo) ConfirmUserEmail(user User) {
	user.ConfirmedAt.Time = time.Now()
	user.ConfirmedAt.Valid = true
	r.db.Save(&user)
}

func (r *UserRepo) ListAllUsers() ([]User, error) {
	var users []User
	r.db.Find(&users)
	return users, nil
}

