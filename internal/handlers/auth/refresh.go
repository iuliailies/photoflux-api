package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleRefresh(ctx *gin.Context) {
	_, ok := common.GetAuthHeader(ctx)
	if !ok {
		return
	}

	var req public.RefreshRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, RefreshError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	// remove token entry from database
	var currtoken model.RefreshToken
	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", req.RefreshToken).Delete(&currtoken).Error

	if err != nil {
		common.EmitError(ctx, RefreshError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not remove current refresh token from database: %s", err.Error())))
		return
	}

	// issue a new access token and a new refresh token
	var tokenstr string
	var refreshstr string

	// The id of the token. Will be associated with the refresh token.
	tokenid := uuid.New()

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   currtoken.UserId.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.accessTokenLifetime)),
		Issuer:    "photoflux",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        tokenid.String(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenstr, err = token.SignedString(h.jwtSecret)
	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not sign auth token: %s", err.Error())))
		return
	}

	// Generate the refresh token.
	refreshToken := model.RefreshToken{
		TokenId: tokenid,
		UserId:  currtoken.UserId,
	}
	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&refreshToken).Error
	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create refresh token: %s", err.Error())))
		return
	}
	refreshstr = refreshToken.Id.String()

	resp := public.RegisterResponse{
		Data: public.AuthData{
			AccessToken:  tokenstr,
			RefreshToken: refreshstr,
		},
	}
	ctx.JSON(http.StatusCreated, &resp)

}
