package stars

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pck/photoflux"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleStarPhoto(ctx *gin.Context) {
	var req public.StarRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, StarPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	Star := model.Star{
		PhotoId: req.PhotoId,
		UserId:  req.UserId,
	}
	isStar := false
	res := h.db.WithContext(ctx).Clauses(clause.Returning{}).Table("stars").
		Where("photo_id = ?", Star.PhotoId).
		Where("user_id = ?", Star.UserId).
		Delete(&Star)

	if err != nil {
		common.EmitError(ctx, StarPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get star data: %s", err.Error())))
		return
	}

	if res.RowsAffected < 1 {
		// Star didn't exist, therefore we create it
		isStar = true
		res = h.db.WithContext(ctx).Clauses(clause.Returning{}).Table("stars").Create(&Star)
	}

	if err != nil {
		common.EmitError(ctx, StarPhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not star the photo: %s", err.Error())))
		return
	}

	resp := public.StarResponse{
		IsStar: isStar,
	}

	ctx.JSON(http.StatusCreated, &resp)
}
