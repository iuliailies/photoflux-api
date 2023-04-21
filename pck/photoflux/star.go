package photoflux

import "github.com/google/uuid"

type StarRequest struct {
	PhotoId uuid.UUID `json:"photo_id"`
	UserId  uuid.UUID `json:"user_id"` // TODO: remove once auth gets implemented
}

type StarResponse struct {
	IsStar bool `json:"is_star"`
}
