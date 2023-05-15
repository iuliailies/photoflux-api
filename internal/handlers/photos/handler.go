package photos

import (
	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/storage"
	"gorm.io/gorm"
)

type handler struct {
	db       *gorm.DB
	apiPaths config.ApiPaths
	storage  *storage.Storage
}

func NewHandler(db *gorm.DB, storage *storage.Storage, config config.Config) handler {
	return handler{
		db:       db,
		apiPaths: config.ApiPaths,
		storage:  storage,
	}
}
