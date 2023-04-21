package model

import (
	"fmt"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	BaseModel
	Name   string  `validate:"required"`
	Photos []Photo `gorm:"many2many:photo_categories;"`
}

// printInfo prints an entity in a simpler format.
func (c *Category) PrintInfo() {
	fmt.Printf("UUID: %s\tNAME: %s\n", c.Id, c.Name)
}

// AfterUpdate is a gorm hook that adds an error if the entry was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty entry.
func (e *Category) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the entry was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty entry.
func (e *Category) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}
