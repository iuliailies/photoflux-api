package photos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleGetPhoto(ctx *gin.Context) {
	_, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	id := ctx.Param("id")
	_uuid, err := uuid.Parse(id)

	if err != nil {
		common.EmitError(ctx, GetPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	var photo model.Photo
	err = h.db.WithContext(ctx).Where("id = ?", _uuid).Take(&photo).Error

	if err != nil {
		common.EmitError(ctx, GetPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get photo: %s", err.Error())))
		return
	}

	var categories []model.Category
	err = h.db.Model(&photo).Association("Categories").Find(&categories)
	if err != nil {
		common.EmitError(ctx, GetPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get photo categories: %s", err.Error())))
		return
	}
	var categoryIds []uuid.UUID = make([]uuid.UUID, 0, len(categories))
	for _, c := range categories {
		categoryIds = append(categoryIds, c.Id)
	}

	starAssociation := h.db.Model(&photo).Association("Users")
	if err != nil {
		common.EmitError(ctx, GetPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get photo star count: %s", err.Error())))
		return
	}

	photoWithStar := model.PhotoWithStars{
		Photo:     photo,
		StarCount: starAssociation.Count(),
	}

	resp := public.GetPhotoResponse{
		Data: PhotoWithRelationshipToPublic(photoWithStar, h.apiPaths, categoryIds),
	}

	ctx.JSON(http.StatusOK, &resp)
}
