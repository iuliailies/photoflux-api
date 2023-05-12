package model

import "github.com/google/uuid"

type RefreshToken struct {
	Id      uuid.UUID `gorm:"default:gen_random_uuid()"`
	TokenId uuid.UUID
	UserId  uuid.UUID
}
