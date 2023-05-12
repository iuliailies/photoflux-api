package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func (h *handler) HandleGetUser(ctx *gin.Context) {
	_, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	id := ctx.Param("id")
	_uuid, err := uuid.Parse(id)

	if err != nil {
		common.EmitError(ctx, GetUserError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	var user model.User
	err = h.db.WithContext(ctx).Where("id = ?", _uuid).Take(&user).Error

	if err != nil {
		common.EmitError(ctx, GetUserError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not get user: %s", err.Error())))
		return
	}

	resp := public.GetUserResponse{
		Data: UserToPublic(user, h.apiPaths, 0),
	}

	ctx.JSON(http.StatusCreated, &resp)

}
