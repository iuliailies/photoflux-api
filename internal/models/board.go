package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardAttr struct {
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	UserId    string    `bson:"user_id"`
	Data      string
}

type Board struct {
	Id        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	UserId    string             `bson:"user_id"`
	Data      string
}

type BoardUpdateAttr struct {
	UpdatedAt time.Time `bson:"updated_at"`
	Data      string
}

// printInfo prints an entity in a simpler format.
func (b *Board) PrintInfo() {
	fmt.Printf("tDATA: %s\n", b.Data)
}
