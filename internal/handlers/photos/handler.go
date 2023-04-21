package photos

import (
	"github.com/iuliailies/photo-flux/internal/config"
	"gorm.io/gorm"
)

type handler struct {
	db       *gorm.DB
	apiPaths config.ApiPaths
}

func NewHandler(db *gorm.DB, config config.Config) handler {
	return handler{
		db:       db,
		apiPaths: config.ApiPaths,
	}
}
