package model

import (
	"fmt"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	BaseModel
	Link       string `validate:"required"`
	IsUploaded bool
	UserId     uuid.UUID  `validate:"required"`
	Categories []Category `gorm:"many2many:photo_categories;"`
	Users      []User     `gorm:"many2many:stars;"`
}

type PhotoWithStars struct {
	Photo
	StarCount int64
}

func (p *Photo) PrintInfo() {
	fmt.Printf("UUID: %s\tNAME: %s\t UPLOAD_TIME: %s\t IS_UPLOADED: %t\n", p.Id, p.Link, p.CreatedAt, p.IsUploaded)
}

// AfterUpdate is a gorm hook that adds an error if the entry was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty entry.
func (e *Photo) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the entry was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty entry.
func (e *Photo) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}
