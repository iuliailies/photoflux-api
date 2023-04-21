package model

import uuid "github.com/google/uuid"

type Star struct {
	PhotoId uuid.UUID
	UserId  uuid.UUID
}
