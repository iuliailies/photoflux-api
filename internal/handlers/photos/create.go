package photos

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	"github.com/iuliailies/photo-flux/internal/rand"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"gorm.io/gorm/clause"
)

type PhotoCategory struct {
	PhotoId    uuid.UUID
	CategoryId uuid.UUID
}

func (h *handler) HandleCreatePhoto(ctx *gin.Context) {
	ah, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	var req public.CreatePhotoRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	if len(req.CategoryIds) == 0 {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not create photo: at least one category should be specified.")))
		return
	}

	// randomly generating a file name in order to uniquily identity it in the minio bucket storage
	objectName, err := rand.RandomStringSecret(64)
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not create photo name: %s", err.Error())))
		return
	}

	Photo := model.Photo{
		UserId: ah.User,
		Name:   objectName,
	}

	objectNameThumbnail := "thumbnail" + objectName
	Thumbnail := Photo
	Thumbnail.Name = objectNameThumbnail

	// ignore thumbnail cause 1 validationis fine
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

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&Thumbnail).Error
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create thumbnail for photo: %s", err.Error())))
		return
	}

	photoCategoryEntries := make([]*PhotoCategory, 0, len(req.CategoryIds))
	for _, cId := range req.CategoryIds {
		entries := []*PhotoCategory{
			{
				PhotoId:    uuid.UUID(Photo.Id),
				CategoryId: uuid.UUID(cId),
			},
			// add the thumbnail too for each category
			{
				PhotoId:    uuid.UUID(Thumbnail.Id),
				CategoryId: uuid.UUID(cId),
			},
		}
		photoCategoryEntries = append(photoCategoryEntries, entries...)
	}

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Table("photo_categories").Create(&photoCategoryEntries).Error

	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create photos and thumbnail associations with categories: %s", err.Error())))
		return
	}

	url_photo, err := h.storage.GetPresignedPut(ctx, "user-"+ah.User.String(), objectName, time.Minute)
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create Photo Href: %s", err.Error())))
		return
	}

	url_thumbnail, err := h.storage.GetPresignedPut(ctx, "user-"+ah.User.String(), objectNameThumbnail, time.Minute)
	if err != nil {
		common.EmitError(ctx, CreatePhotoError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create Photo Href: %s", err.Error())))
		return
	}

	photoWithStar := model.PhotoWithStars{
		Photo:     Photo,
		StarCount: 0,
	}

	resp := public.CreatePhotoResponse{
		Data: PhotoWithRelationshipToPublic(photoWithStar, h.apiPaths, req.CategoryIds, url_photo, url_thumbnail),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
