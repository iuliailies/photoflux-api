package main

import (
	"fmt"
	"os"

	log "github.com/Ozoniuss/stdlog"
	"github.com/iuliailies/photo-flux/internal/config"
	pfrabbit "github.com/iuliailies/photo-flux/internal/rabbitmq"
	"github.com/iuliailies/photo-flux/internal/router"
	"github.com/iuliailies/photo-flux/internal/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connect starts the connection with the database.
func connect(config config.Database) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Name),
	))
	return conn, err
}

func run() error {

	c, err := config.ParseConfig()
	if err != nil {
		return fmt.Errorf("could not initialize config: %w", err)
	}
	fmt.Printf("config: %+v\n", c)

	db, err := connect(c.Database)

	if err != nil {
		return fmt.Errorf("could not connect to db: %w", err)
	}

	storage, err := storage.New(c.Storage)
	if err != nil {
		return fmt.Errorf("could not initialize minio connection: %w", err)
	}

	// Start the notification listener
	uploadsListener := pfrabbit.NewUploadsListener(db, c.Notifications.RabbitMQ)
	err = uploadsListener.Start()
	if err != nil {
		return fmt.Errorf("could not start uploads listener: %w", err)
	}

	engine, err := router.NewRouter(db, storage, c)
	if err != nil {
		return fmt.Errorf("could not initialize router: %w", err)
	}
	engine.Run(":8033")

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Errf("Error running api: %s", err.Error())
		os.Exit(1)
	}
}
