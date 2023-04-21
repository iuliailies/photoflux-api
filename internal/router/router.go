package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/endpoints"
	"github.com/iuliailies/photo-flux/internal/handlers/categories"
	"github.com/iuliailies/photo-flux/internal/handlers/photos"
	"github.com/iuliailies/photo-flux/internal/handlers/stars"
	"github.com/iuliailies/photo-flux/internal/handlers/users"
	"gorm.io/gorm"
)

// NewRouter initializes the gin router with the existing handlers and options.
func NewRouter(db *gorm.DB, config config.Config) (*gin.Engine, error) {
	r := gin.Default()
	{
		h := photos.NewHandler(db, config)
		endpoints.RegisterPhotos(r, &h)
	}
	{
		h := users.NewHandler(db, config)
		endpoints.RegisterUsers(r, &h)
	}
	{
		h := categories.NewHandler(db, config)
		endpoints.RegisterCategories(r, &h)
	}
	{
		h := stars.NewHandler(db, config)
		endpoints.RegisterStars(r, &h)
	}
	return r, nil
}
