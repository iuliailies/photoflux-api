package boards

import (
	"github.com/iuliailies/photo-flux/internal/config"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func BoardToItem(board model.Board, apipath config.ApiPaths) public.BoardData {
	return public.BoardData{
		ResourceID: public.ResourceID{
			Id:   board.Id.Hex(),
			Type: public.BoardType,
		},
		UserId: board.UserId,
		Timestamps: public.Timestamps{
			CreatedAt: board.CreatedAt,
			UpdatedAt: board.UpdatedAt,
		},
		Data: board.Data,
	}
}
