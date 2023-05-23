package photos

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/config"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func PhotoToPublic(photo model.PhotoWithStars, apipath config.ApiPaths, url string) public.PhotoData {
	return public.PhotoData{
		ResourceID: public.ResourceID{
			Id:   photo.Id.String(),
			Type: public.PhotoType,
		},
		Attributes: public.PhotoAttributes{
			IsUploaded: false,
			UserId:     photo.UserId,
			Timestamps: public.Timestamps{
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			},
		},
		Meta: public.PhotoMeta{
			NumberStars: photo.StarCount,
			HRef:        url,
		},
		Links: public.PhotoLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Photos, photo.Id.String()),
		},
	}
}

func PhotoWithRelationshipToPublic(photo model.PhotoWithStars, apipath config.ApiPaths, categoryIds []uuid.UUID, url string) public.PhotoData {
	photoCategoryEntries := make([]public.PhotoCategoryData, 0, len(categoryIds))
	for _, cId := range categoryIds {
		entry := public.PhotoCategoryData{
			Type: public.CategoryType,
			Id:   uuid.UUID(cId),
		}
		photoCategoryEntries = append(photoCategoryEntries, entry)
	}
	return public.PhotoData{
		ResourceID: public.ResourceID{
			Id:   photo.Id.String(),
			Type: public.PhotoType,
		},
		Attributes: public.PhotoAttributes{
			IsUploaded: false,
			UserId:     photo.UserId,
			Timestamps: public.Timestamps{
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			},
		},
		Meta: public.PhotoMeta{
			NumberStars: photo.StarCount,
			HRef:        url,
		},
		Links: public.PhotoLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Photos, photo.Id.String()),
		},
		Relationships: public.PhotoRelationships{
			Categories: public.PhotoCategoriesRelationship{
				Links: public.CategoryLinks{
					Self: apipath.Categories,
				},
				Data: photoCategoryEntries,
			},
		},
	}
}

func PhotoToPublicListItem(photo model.PhotoWithStars, apipath config.ApiPaths, url string, starred bool) public.PhotoListItemData {
	return public.PhotoListItemData{
		ResourceID: public.ResourceID{
			Id:   photo.Id.String(),
			Type: public.PhotoType,
		},
		Attributes: public.PhotoAttributes{
			IsUploaded: false,
			UserId:     photo.UserId,
			Timestamps: public.Timestamps{
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			},
		},
		Meta: public.PhotoMeta{
			NumberStars:   photo.StarCount,
			HRef:          url,
			StarredByUser: starred,
		},
		Links: public.PhotoListItemLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Photos, photo.Id.String()),
		},
	}
}
