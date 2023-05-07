package photos

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/config"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func PhotoToPublic(photo model.PhotoWithStars, apipath config.ApiPaths) public.PhotoData {
	return public.PhotoData{
		ResourceID: public.ResourceID{
			Id:   photo.Id.String(),
			Type: public.PhotoType,
		},
		Attributes: public.PhotoAttributes{
			Link:       photo.Link,
			IsUploaded: false,
			UserId:     photo.UserId,
			Timestamps: public.Timestamps{
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			},
		},
		Meta: public.PhotoMeta{
			NumberStars: photo.StarCount,
		},
		Links: public.PhotoLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Photos, photo.Id.String()),
		},
	}
}

func PhotoWithRelationshipToPublic(photo model.PhotoWithStars, apipath config.ApiPaths, categoryIds []uuid.UUID) public.PhotoData {
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
			Link:       photo.Link,
			IsUploaded: false,
			UserId:     photo.UserId,
			Timestamps: public.Timestamps{
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			},
		},
		Meta: public.PhotoMeta{
			NumberStars: photo.StarCount,
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

func PhotoToPublicListItem(photo model.PhotoWithStars, apipath config.ApiPaths) public.PhotoListItemData {
	return public.PhotoListItemData{
		ResourceID: public.ResourceID{
			Id:   photo.Id.String(),
			Type: public.PhotoType,
		},
		Attributes: public.PhotoAttributes{
			Link:       photo.Link,
			IsUploaded: false,
			UserId:     photo.UserId,
			Timestamps: public.Timestamps{
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			},
		},
		Meta: public.PhotoMeta{
			NumberStars: photo.StarCount,
		},
		Links: public.PhotoListItemLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Photos, photo.Id.String()),
		},
	}
}
