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

	categoryIds := public.CategoriesFromURL(*params.Category)
	var categories []model.Category
	err = h.db.WithContext(ctx).Where("id IN (?)", categoryIds).Find(&categories).Error

	if err != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list photos. Invalid category: %s", err.Error())))
		return
	}

	names := make([]string, len(categories))
	for i, category := range categories {
		names[i] = category.Name
	}

	sortCriteria := "created_at"
	if params.Sort != nil {
		if *params.Sort == "star" {
			sortCriteria = "star_count"
		} else if *params.Sort != "created_at" {
			common.EmitError(ctx, ListPhotoError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not list photos. The `sort` parameter should either be `star` or `created_at`.")))
			return
		}
	}

	cols := []gorm.OrderedColumn{
		{
			Column: sortCriteria,
			Order:  gorm.OrderDESC,
		},
		{
			Column: "id",
			Order:  gorm.OrderASC,
		},
	}

	var photos []model.PhotoWithStars

	selection := h.db.Debug().Table("photos").
		Joins("JOIN photo_categories ON photo_categories.photo_id = photos.id").
		Joins("LEFT JOIN stars ON stars.photo_id = photos.id").
		Where("photo_categories.category_id IN (?)", categoryIds).
		Group("photos.id").
		Having("COUNT(DISTINCT photo_categories.category_id) >= ?", len(categoryIds)-1).
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
			CategoryName: names,
		},
		Links: public.ListPhotoLinks{
			Next: model.BuildNextLink(nextarr, "photos/", params.Limit, params.Category),
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
