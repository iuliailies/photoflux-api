package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id        uuid.UUID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	UserId    string    `bson:"user_id"`
	Data      string    `bson:"data"`
}

// printInfo prints an entity in a simpler format.
func (b *Board) PrintInfo() {
	fmt.Printf("UUID: %s\tDATA: %s\n", b.Id, b.Data)
}
