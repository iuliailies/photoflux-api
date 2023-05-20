package postgres

import (
	"fmt"

	"github.com/iuliailies/photo-flux/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config config.Database) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Name),
	))
	return conn, err
}
