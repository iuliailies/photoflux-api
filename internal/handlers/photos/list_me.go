package photos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleListMyPhoto(ctx *gin.Context) {
	ah, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	var params public.ListMyPhotoParams
	err := ctx.ShouldBindQuery(&params)

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	filters["photos.user_id"] = ah.User

	var photos []model.PhotoWithStars
	var count int64

	err = h.db.Debug().Table("photos").
		Joins("JOIN photo_categories ON photo_categories.photo_id = photos.id").
		Joins("LEFT JOIN stars ON stars.photo_id = photos.id").
		Where(filters).
		Group("photos.id").
		Select("photos.id, photos.link, photos.user_id, photos.is_uploaded, photos.created_at, photos.updated_at, COUNT(stars.user_id) AS star_count").
		Order("created_at DESC").
		Scan(&photos).
		Count(&count).
		Error

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list user photos: %s", err.Error())))
		return
	}

	var totalStars int64
	err = h.db.Debug().
		Table("photos").
		Select("COUNT(*) as total_stars").
		Joins("JOIN stars ON photos.id = stars.photo_id").
		Where("photos.user_id = ?", ah.User.ID()).
		Scan(&totalStars).
		Error

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get total star count of user: %s", err.Error())))
		return
	}

	resp := public.ListPhotoResponse{
		Data: make([]public.PhotoListItemData, 0, len(photos)),
		Meta: public.PhotoListMeta{
			NumberStars:  totalStars,
			NumberPhotos: count,
		},
		Links: public.ListPhotoLinks{
			Self: h.apiPaths.Photos + "/",
		},
	}
	for _, photo := range photos {
		resp.Data = append(resp.Data, PhotoToPublicListItem(photo, h.apiPaths))
	}
	ctx.JSON(http.StatusOK, &resp)
}
