package auth

import (
	"time"

	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/storage"
	"gorm.io/gorm"
)

type handler struct {
	db                  *gorm.DB
	storage             *storage.Storage
	apiPaths            config.ApiPaths
	jwtSecret           []byte
	accessTokenLifetime time.Duration
	minioTokenLifetime  time.Duration
}

func NewHandler(db *gorm.DB, storage *storage.Storage, config config.Config) handler {

	return handler{
		db:                  db,
		apiPaths:            config.ApiPaths,
		jwtSecret:           config.Auth.Secret,
		storage:             storage,
		accessTokenLifetime: config.Auth.AccessTokenLifetime,
		minioTokenLifetime:  config.Auth.MinioTokenLifetime,
	}
}
