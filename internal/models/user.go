package model

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name     string  `validate:"required"`
	Email    string  `validate:"required"`
	Password string  `validate:"required"`
	Photos   []Photo `gorm:"many2many:stars;"`
}

func (u *User) PrintInfo() {
	fmt.Printf("UUID: %s\tNAME: %s\t EMAIL: %s\t PASSWORD: %s\n", u.Id, u.Name, u.Email, u.Password)
}

// AfterUpdate is a gorm hook that adds an error if the entry was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty entry.
func (e *User) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the entry was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty entry.
func (e *User) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

// BeforeUpdate : hook before a user is updated
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	fmt.Println("before update")
	fmt.Println(u.Password)

	if u.Password != "" {
		hash, err := MakePassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
		// tx.Model(u).Update("password", hash)
	}

	return
}

// MakePassword : Encrypt user password
func MakePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
