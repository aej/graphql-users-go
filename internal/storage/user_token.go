package storage

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type UsersToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	User      User
	Token     []byte         `gorm:"type:bytes;not null"`
	Context   string         `gorm:"type:varchar(16)"`
	CreatedAt time.Time      `gorm:"not null"`
	SentTo    sql.NullString `gorm:"type:varchar(255);null"`
}

func (r *UserRepo) CreateUserToken(UserId uuid.UUID, token []byte, context string, sentTo sql.NullString) (UsersToken, error) {
	userToken := UsersToken{UserId: UserId, Token: token, Context: context, SentTo: sentTo}
	r.db.Create(&userToken)

	return userToken, nil
}

// Delete All UsersTokens for an associated User and context
func (r *UserRepo) DeleteUserToken(UserId uuid.UUID, context string) {
	var userTokens []UsersToken
	result := r.db.Where("user_id = ? AND context = ?", UserId, context).Find(&userTokens)

	if result.RowsAffected > 0 {
		var idsToDelete []string
		for _, token := range userTokens {
			idsToDelete = append(idsToDelete, token.ID.String())
		}
		r.db.Delete(&userTokens, idsToDelete) // delete all tokens with this primary key
	}
}

// This function needs to preload the user object
func (r *UserRepo) GetUnexpiredUserTokenForContext(token []byte, context string, validUntil time.Time) (UsersToken, error) {
	var usersToken UsersToken

	result := r.db.Preload("User").Where("token = ? AND context = ? AND created_at > ?", token, context, validUntil).First(&usersToken)

	if result.RowsAffected == 1 {
		return usersToken, nil
	}
	return usersToken, RecordNotFound
}
