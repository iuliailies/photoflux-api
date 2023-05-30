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
	"gorm.io/gorm/clause"
)

func (h *handler) HandleListPhoto(ctx *gin.Context) {
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

	var params public.ListPhotoParams
	err := ctx.ShouldBindQuery(&params)

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusBadRequest,
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

	var category model.Category
	err = h.db.WithContext(ctx).Where("id = ?", *params.Category).Take(&category).Error

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list photos. Invalid category: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	if params.Category != nil {
		filters["photo_categories.category_id"] = *params.Category
	}
	// filters["photos.is_uploaded"] = true // TODO: recheck bug

	var photos []model.PhotoWithStars

	selection := h.db.Debug().Table("photos").
		Joins("JOIN photo_categories ON photo_categories.photo_id = photos.id").
		Joins("LEFT JOIN stars ON stars.photo_id = photos.id").
		Where(filters).
		Group("photos.id").
		Select("photos.id, photos.user_id, photos.name, photos.is_uploaded, photos.created_at, photos.updated_at, COUNT(stars.user_id) AS star_count")

	photos, _, nextarr, errarr := gorm.ListMultiColumn(
		selection,
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
			fmt.Sprintf("Could not list photos: %s", err.Error())))
		return
	}

	resp := public.ListPhotoResponse{
		Data: make([]public.PhotoListItemData, 0, len(photos)),
		Meta: public.PhotoListMeta{
			CategoryName: category.Name,
		},
		Links: public.ListPhotoLinks{
			Next: model.BuildNextLink(nextarr, "photos/", params.Limit),
		},
	}
	for _, photo := range photos {
		var star []model.Star
		err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Table("stars").
			Where("photo_id = ?", photo.Id.String()).
			Where("user_id = ?", ah.User).
			Find(&star).
			Error

		if err != nil {
			common.EmitError(ctx, ListPhotoError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not get star data: %s", err.Error())))
			return
		}

		isStarredByUser := len(star) == 1

		// TODO error handling
		url, _ := h.storage.GetPresignedGet(ctx, "user-"+photo.UserId.String(), photo.Name, time.Minute)

		resp.Data = append(resp.Data, PhotoToPublicListItem(photo, h.apiPaths, url, isStarredByUser))
	}
	ctx.JSON(http.StatusOK, &resp)
}
