package photoflux

import "github.com/google/uuid"

type StarRequest struct {
	PhotoId uuid.UUID `json:"photo_id"`
}

type StarResponse struct {
	IsStar bool `json:"is_star"`
}
