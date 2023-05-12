package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleLogout(ctx *gin.Context) {
	var req public.LogoutRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	var currtoken model.RefreshToken
	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", req.RefreshToken).Delete(&currtoken).Error

	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not delete refresh token: %s", err.Error())))
		return
	}

	resp := public.LogoutResponse{}
	ctx.JSON(http.StatusCreated, &resp)

}
