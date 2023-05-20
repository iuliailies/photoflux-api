package boards

import (
	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	db       *gorm.DB
	mongoDb  *mongo.Client
	apiPaths config.ApiPaths
	storage  *storage.Storage
}

func NewHandler(db *gorm.DB, mdb *mongo.Client, storage *storage.Storage, config config.Config) handler {
	return handler{
		db:       db,
		mongoDb:  mdb,
		apiPaths: config.ApiPaths,
		storage:  storage,
	}
}
