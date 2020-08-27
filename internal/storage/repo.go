package storage

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Repositories struct {
	User UserRepo
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	repositories := &Repositories{
		User: UserRepo{db: db},
	}
	return repositories, nil
}

var RecordNotFound = errors.New("record could not be found")
