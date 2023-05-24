package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardAttr struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    string
	Data      string
}

type Board struct {
	Id primitive.ObjectID
	BoardAttr
}

// printInfo prints an entity in a simpler format.
func (b *Board) PrintInfo() {
	fmt.Printf("tDATA: %s\n", b.Data)
}
