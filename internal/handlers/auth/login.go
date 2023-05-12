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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleLogin(ctx *gin.Context) {
	var req public.LoginRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	filters["email"] = req.Email

	// check that the user's details (email, pass etc) are valid
	var user model.User
	err = h.db.WithContext(ctx).First(&user, filters).Clauses(clause.Returning{}).Error

	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not login user. Invalid email: %s", err.Error())))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		common.EmitError(ctx, LoginError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not login user. Invalid password: %s", err.Error())))
		return
	}

	// create a new access token, and a new refresh token
	var tokenstr string
	var refreshstr string

	// The id of the token. Will be associated with the refresh token.
	tokenid := uuid.New()

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.Id.String(),
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
		UserId:  user.Id,
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
