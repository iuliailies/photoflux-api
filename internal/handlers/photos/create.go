package photos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pck/photoflux"
	"gorm.io/gorm/clause"
)

type PhotoCategory struct {
	PhotoId    uuid.UUID
	CategoryId uuid.UUID
}

func (h *handler) HandleCreatePhoto(ctx *gin.Context) {
	var req public.CreatePhotoRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	if len(req.CategoryIds) == 0 {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not create photo: at least one category should be spcified.")))
		return
	}

	Photo := model.Photo{
		Link:   req.Link,
		UserId: req.UserId,
	}

	v := validator.New()
	err = v.Struct(Photo)
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create photo: %s", err.Error())))
		return
	}

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&Photo).Error
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create photo: %s", err.Error())))
		return
	}

	photoCategoryEntries := make([]*PhotoCategory, 0, len(req.CategoryIds))
	for _, cId := range req.CategoryIds {
		entry := PhotoCategory{
			PhotoId:    uuid.UUID(Photo.Id),
			CategoryId: uuid.UUID(cId),
		}
		photoCategoryEntries = append(photoCategoryEntries, &entry)
	}

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Table("photo_categories").Create(&photoCategoryEntries).Error

	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create Photo: %s", err.Error())))
		return
	}

	photoWithStar := model.PhotoWithStars{
		Photo:     Photo,
		StarCount: 0,
	}

	resp := public.CreatePhotoResponse{
		Data: PhotoWithRelationshipToPublic(photoWithStar, h.apiPaths, req.CategoryIds),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
