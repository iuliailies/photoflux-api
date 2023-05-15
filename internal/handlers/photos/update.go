package photos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleUpdatePhoto(ctx *gin.Context) {
	_, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)

	if err != nil {
		common.EmitError(ctx, UpdatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	var req public.UpdatePhotoRequest
	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, GetPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make(map[string]any)

	if req.IsUploaded != nil {
		updatedFields["is_uploaded"] = *req.IsUploaded
	}

	var photo model.Photo
	err = h.db.WithContext(ctx).Model(&photo).Clauses(clause.Returning{}).Where("id = ?", uuid).Updates(updatedFields).Error

	if err != nil {
		common.EmitError(ctx, UpdatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not update photo: %s", err.Error())))
		return
	}

	starAssociation := h.db.Model(&photo).Association("Users")
	if err != nil {
		common.EmitError(ctx, UpdatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get photo star number: %s", err.Error())))
		return
	}

	photoWithStars := model.PhotoWithStars{
		Photo:     photo,
		StarCount: starAssociation.Count(),
	}

	// TODO: integrate with minio if needed
	resp := public.UpdatePhotoResponse{
		Data: PhotoToPublic(photoWithStars, h.apiPaths, ""),
	}

	ctx.JSON(http.StatusOK, &resp)
}
