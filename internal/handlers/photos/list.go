package photos

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleListPhoto(ctx *gin.Context) {
	_, ok := common.GetAuthHeader(ctx)
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

	var photos []model.PhotoWithStars

	selection := h.db.Debug().Table("photos").
		Joins("JOIN photo_categories ON photo_categories.photo_id = photos.id").
		Joins("LEFT JOIN stars ON stars.photo_id = photos.id").
		Where(filters).
		Group("photos.id").
		Select("photos.id, photos.user_id, photos.name, photos.is_uploaded, photos.created_at, photos.updated_at, COUNT(stars.user_id) AS star_count").
		Scan(&photos)

	if params.Sort != nil {
		fmt.Println("inside params sort", *params.Sort)
		var orderString string
		if *params.Sort == "created_at" {
			orderString = "created_at DESC"
		} else if *params.Sort == "star" {
			orderString = "star_count DESC"
		} else {
			common.EmitError(ctx, ListPhotoError(
				http.StatusBadRequest,
				`Could not list photos: "sort" query parameter should be either "star" or "created_at".`))
			return
		}
		selection = selection.Order(orderString).Scan(&photos)
	}

	if selection.Error != nil {
		common.EmitError(ctx, ListPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list photos: %s", selection.Error.Error())))
		return
	}

	resp := public.ListPhotoResponse{
		Data: make([]public.PhotoListItemData, 0, len(photos)),
		Meta: public.PhotoListMeta{
			CategoryName: category.Name,
		},
		Links: public.ListPhotoLinks{
			Self: h.apiPaths.Photos + "/",
		},
	}
	for _, photo := range photos {
		// TODO error handling
		url, _ := h.storage.GetPresignedGet(ctx, "user-"+photo.UserId.String(), photo.Name, time.Minute)
		resp.Data = append(resp.Data, PhotoToPublicListItem(photo, h.apiPaths, url))
	}
	ctx.JSON(http.StatusOK, &resp)
}
