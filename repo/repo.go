package repo

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID uuid.UUID		  `gorm:"type:uuid;unique_index;default:generate_uuid_v4()"`
	Email string 		  `gorm:"type:varchar(255);unique_index;not_null"`
	HashedPassword string `gorm:"type:not_null"`
	CreatedAt time.Time   `gorm:"type:not_null"`
	UpdatedAt time.Time	  `gorm:"type:not_null"`
}

type UserRepo struct {
	db *gorm.DB
}

func(r *UserRepo) CreateUser(email string) (User, error) {
	user := User{Email: email}
	r.db.Create(user)
	return user, nil
}

func(r *UserRepo) ListAllUsers() ([]User, error) {
	var users []User
	r.db.Find(&users)

	return users, nil
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