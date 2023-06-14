package photos

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/gorm"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleListMyPhoto(ctx *gin.Context) {
	ah, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	cols := []gorm.OrderedColumn{
		{
			Column: "created_at",
			Order:  gorm.OrderDESC,
		},
		{
			Column: "id",
			Order:  gorm.OrderASC,
		},
	}

	var params public.ListMyPhotoParams
	err := ctx.ShouldBindQuery(&params)

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	var afterarr []any = nil

	if params.After != nil {
		var photoCursor model.PhotoCursor = model.FromURLString(*params.After)
		afterarr = []any{photoCursor.CreatedAt, photoCursor.Id}
	}

	limit := -1

	if params.Limit != nil {
		limit = *params.Limit
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	filters["photos.user_id"] = ah.User
	// filters["photos.is_uploaded"] = true

	var photos []model.PhotoWithStars

	tdb := h.db.Debug().Table("photos").
		Joins("LEFT JOIN stars ON stars.photo_id = photos.id").
		Where(filters).
		Group("photos.id").
		Select("photos.id, photos.user_id, photos.name, photos.is_uploaded, photos.created_at, photos.updated_at, COUNT(stars.user_id) AS star_count")

	photos, _, nextarr, errarr := gorm.ListMultiColumn(
		tdb,
		"photos",
		limit,
		cols,
		model.RetrieveCursorArr,
		nil,
		afterarr,
		gorm.OrderASC,
	)

	if errarr != nil {
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
		Where("photos.user_id = ?", ah.User).
		Scan(&totalStars).
		Error

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get total star count for user: %s", err.Error())))
		return
	}

	var count int64
	err = h.db.Debug().
		Table("photos").
		Select("COUNT(*) as total_photos").
		Where("photos.user_id = ?", ah.User).
		Scan(&count).
		Error

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get total photo number for user: %s", err.Error())))
		return
	}

	resp := public.ListMyPhotoResponse{
		Data: make([]public.PhotoListItemData, 0, len(photos)),
		Meta: public.MyPhotoListMeta{
			NumberStars:  totalStars,
			NumberPhotos: count,
		},
		Links: public.ListPhotoLinks{
			Next: model.BuildNextLink(nextarr, "photos/me/", params.Limit, nil),
		},
	}
	for _, photo := range photos {
		// TODO error handling
		url, _ := h.storage.GetPresignedGet(ctx, "user-"+ah.User.String(), photo.Name, time.Minute)
		resp.Data = append(resp.Data, PhotoToPublicListItem(photo, h.apiPaths, url, false))
	}
	ctx.JSON(http.StatusOK, &resp)
}
