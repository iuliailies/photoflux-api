package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pck/photoflux"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleCreateUser(ctx *gin.Context) {
	var req public.CreateUserRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, CreateUserError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	User := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	v := validator.New()
	err = v.Struct(User)
	if err != nil {
		common.EmitError(ctx, CreateUserError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create user: %s", err.Error())))
		return
	}

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&User).Error

	if err != nil {
		common.EmitError(ctx, CreateUserError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create user: %s", err.Error())))
		return
	}

	resp := public.CreateUserResponse{
		Data: UserToPublic(User, h.apiPaths, 0),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
