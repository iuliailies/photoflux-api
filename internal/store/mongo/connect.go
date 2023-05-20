package mongo

import (
	"context"
	"fmt"

	"github.com/iuliailies/photo-flux/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, config config.MongoDatabase) (*mongo.Client, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)), &options.ClientOptions{
		Auth: &options.Credential{
			Username: config.User,
			Password: config.Password,
		},
	})
	return client, err
}
