package users

import (
	"fmt"

	"github.com/iuliailies/photo-flux/internal/config"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func UserToPublic(user model.User, apipath config.ApiPaths, runningTotal float32) public.UserData {
	return public.UserData{
		ResourceID: public.ResourceID{
			Id:   user.Id.String(),
			Type: public.UserType,
		},
		Attributes: public.UserAttributes{
			Name:  user.Name,
			Email: user.Email,
			Timestamps: public.Timestamps{
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
		Links: public.UserLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Users, user.Id.String()),
		},
	}
}
