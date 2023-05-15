package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/endpoints"
	"github.com/iuliailies/photo-flux/internal/handlers/auth"
	"github.com/iuliailies/photo-flux/internal/handlers/categories"
	"github.com/iuliailies/photo-flux/internal/handlers/photos"
	"github.com/iuliailies/photo-flux/internal/handlers/stars"
	"github.com/iuliailies/photo-flux/internal/handlers/users"
	"github.com/iuliailies/photo-flux/internal/storage"
	"gorm.io/gorm"
)

// NewRouter initializes the gin router with the existing handlers and options.
func NewRouter(db *gorm.DB, storage *storage.Storage, config config.Config) (*gin.Engine, error) {
	r := gin.Default()
	SetupCORS(r)
	{
		h := photos.NewHandler(db, storage, config)
		endpoints.RegisterPhotos(r, config.Auth, &h)
	}
	{
		h := users.NewHandler(db, storage, config)
		endpoints.RegisterUsers(r, config.Auth, &h)
	}
	{
		h := categories.NewHandler(db, config)
		endpoints.RegisterCategories(r, config.Auth, &h)
	}
	{
		h := stars.NewHandler(db, config)
		endpoints.RegisterStars(r, config.Auth, &h)
	}
	{
		h := auth.NewHandler(db, storage, config)
		endpoints.RegisterAuth(r, config.Auth, &h)
	}
	return r, nil
}
