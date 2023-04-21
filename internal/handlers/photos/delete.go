package photos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pck/photoflux"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeletePhoto(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	var photo model.Photo
	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", uuid).Delete(&photo).Error

	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not delete photo: %s", err.Error())))
		return
	}

	photoWithStars := model.PhotoWithStars{
		Photo:     photo,
		StarCount: 0,
	}

	resp := public.DeletePhotoResponse{
		Data: PhotoToPublic(photoWithStars, h.apiPaths),
	}

	resp.Data.Links.Self = ""

	ctx.JSON(http.StatusOK, &resp)
}
