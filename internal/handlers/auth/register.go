package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Ozoniuss/stdlog"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/handlers/common"
	model "github.com/iuliailies/photo-flux/internal/models"
	pfrand "github.com/iuliailies/photo-flux/internal/rand"
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleRegister(ctx *gin.Context) {

	var req public.RegisterRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, RegisterError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	v := validator.New()
	err = v.Struct(user)
	if err != nil {
		common.EmitError(ctx, RegisterError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid user data: %s", err.Error())))
		return
	}

	var tokenstr string
	var refreshstr string

	// All error handling is done within the transaction.
	err = h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txerr := tx.Clauses(clause.Returning{}).Create(&user).Error
		if txerr != nil {
			return RegisterError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not create user: %s", txerr.Error()))
		}

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
		tokenstr, txerr = token.SignedString(h.jwtSecret)
		if txerr != nil {
			return RegisterError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not sign auth token: %s", txerr.Error()),
			)
		}

		// Generate the refresh token.
		refreshToken := model.RefreshToken{
			TokenId: tokenid,
			UserId:  user.Id,
		}
		txerr = tx.Clauses(clause.Returning{}).Create(&refreshToken).Error
		if txerr != nil {
			return RegisterError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not create refresh token: %s", txerr.Error()))
		}
		refreshstr = refreshToken.Id.String()

		// we need to generate a secure password befor creating a minio user
		// however, this will not be stored, nor reused later
		userSecret, err := pfrand.RandomStringSecret(32)
		if err != nil {
			panic(err)
		}

		txerr = h.storage.NewMinioUser(ctx, user.Id, userSecret)
		if txerr != nil {
			stdlog.Errf("Could not create minio user: %s\n", txerr.Error())
			return RegisterError(
				http.StatusInternalServerError,
				"An error occured during the creation of the storage account.",
			)
		}
		return nil
	})

	if err != nil {
		var perr public.Error
		switch {
		case errors.As(err, &perr):
			common.EmitError(ctx, perr)
			return
		default:
			common.EmitError(ctx, RegisterError(
				http.StatusInternalServerError,
				fmt.Sprintf("An unknown occured during user registration: %s", err.Error())))
			return
		}
	}

	resp := public.RegisterResponse{
		Data: public.AuthData{
			AccessToken:  tokenstr,
			RefreshToken: refreshstr,
		},
	}
	ctx.JSON(http.StatusCreated, &resp)
}
