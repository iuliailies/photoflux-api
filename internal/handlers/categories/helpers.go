package categories

import (
	"fmt"

	"github.com/iuliailies/photo-flux/internal/config"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pck/photoflux"
)

func CategoryToPublicListItem(category model.Category, apipath config.ApiPaths) public.CategoryListItemData {
	return public.CategoryListItemData{
		ResourceID: public.ResourceID{
			Id:   category.Id.String(),
			Type: public.PhotoType,
		},
		Attributes: public.CategoryAttributes{
			Name: category.Name,
			Timestamps: public.Timestamps{
				CreatedAt: category.CreatedAt,
				UpdatedAt: category.UpdatedAt,
			},
		},
		Links: public.CategoryListItemLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Categories, category.Id.String()),
		},
	}
}
