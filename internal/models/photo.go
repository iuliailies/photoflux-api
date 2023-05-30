package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	BaseModel
	IsUploaded bool
	UserId     uuid.UUID `validate:"required"`
	Name       string
	Categories []Category `gorm:"many2many:photo_categories;"`
	Users      []User     `gorm:"many2many:stars;"`
}

type PhotoCursor struct {
	Id        uuid.UUID
	CreatedAt time.Time
}

type PhotoWithStars struct {
	Photo
	StarCount int64
}

func (p *Photo) PrintInfo() {
	fmt.Printf("UUID: %s\t UPLOAD_TIME: %s\t IS_UPLOADED: %t\n", p.Id, p.CreatedAt, p.IsUploaded)
}

func RetrieveCursorArr(p *PhotoWithStars) []any {
	return []any{p.CreatedAt, p.Id}
}

func (m PhotoCursor) toURLString() string {
	return strings.Join([]string{m.Id.String(), strconv.FormatInt(m.CreatedAt.UnixNano(), 10)}, ".")
}

func FromURLString(urlstr string) PhotoCursor {
	values := strings.Split(urlstr, ".")

	tm, _ := strconv.ParseInt(values[1], 10, 64)

	m := PhotoCursor{
		Id:        uuid.MustParse(values[0]), //TODO: error handling
		CreatedAt: time.Unix(0, tm),
	}

	return m
}

func BuildNextLink(arr []any, preString string, limit *int) string {
	if len(arr) == 0 {
		return ""
	}
	var cursor PhotoCursor = PhotoCursor{
		Id:        arr[1].(uuid.UUID),
		CreatedAt: arr[0].(time.Time),
	}
	str := cursor.toURLString()
	link := fmt.Sprintf("%s?after=%s", preString, str)
	if limit != nil {
		link = link + fmt.Sprintf("\u0026limit=%d", *limit)
	}
	return link
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
