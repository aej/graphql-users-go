package userservice

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	Email          string    `gorm:"type:varchar(255);unique_index;not null"`
	HashedPassword string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

type UsersToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	Token     []byte    `gorm:"type:bytes;not null"`
	CreatedAt time.Time `gorm:"not null"`
}
