package main

import (
	"fmt"
	"os"

	log "github.com/Ozoniuss/stdlog"
	"github.com/iuliailies/photo-flux/internal/config"
	"github.com/iuliailies/photo-flux/internal/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connect starts the connection with the database.
func connect() (*gorm.DB, error) {
	host := "localhost"
	port := 5432
	user := "photoflux"
	dbname := "photoflux"
	password := "photoflux"

	conn, err := gorm.Open(postgres.Open(
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, dbname),
	))
	return conn, err
}

func run() error {
	db, err := connect()

	if err != nil {
		return fmt.Errorf("could not connect to db: %w", err)
	}

	// TODO: actually use config

	c := config.Config{}

	engine, err := router.NewRouter(db, c)
	engine.Run("127.0.0.1:8033")

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Errf("Error running api: %s", err.Error())
		os.Exit(1)
	}
}
