package photoflux

import "github.com/google/uuid"

const PhotoType = "photo"

// Data returned about a photo when a single one is returned.
type PhotoData struct {
	ResourceID
	Attributes    PhotoAttributes    `json:"attributes"`
	Meta          PhotoMeta          `json:"meta"`
	Links         PhotoLinks         `json:"links"`
	Relationships PhotoRelationships `json:"relationships"`
}

// Data returned about a photo when a list of photos is returned
type PhotoListItemData struct {
	ResourceID
	Attributes PhotoAttributes    `json:"attributes"`
	Meta       PhotoMeta          `json:"meta"`
	Links      PhotoListItemLinks `json:"links"`
}

type PhotoAttributes struct {
	IsUploaded bool      `json:"is_uploaded"`
	UserId     uuid.UUID `json:"user_id"`
	Timestamps
}

type PhotoMeta struct {
	NumberStars   int64  `json:"number_stars"`
	HRef          string `json:"href"`
	HRefThumbnail string `json:"href_thumbnail"`
	StarredByUser bool   `json:"starred_by_user"`
}

type PhotoListMeta struct {
	CategoryName []string `json:"category_name"`
}

type MyPhotoListMeta struct {
	NumberStars  int64 `json:"number_stars"`
	NumberPhotos int64 `json:"number_photos"`
}

type PhotoLinks struct {
	Self string `json:"self"`
}

type PhotoListItemLinks struct {
	Self string `json:"self"`
}

type PhotoRelationships struct {
	Categories PhotoCategoriesRelationship `json:"categories"`
}

type PhotoCategoriesRelationship struct {
	Links CategoryLinks       `json:"links"`
	Data  []PhotoCategoryData `json:"data"`
}

type PhotoCategoryData struct {
	Type string    `json:"type"`
	Id   uuid.UUID `json:"id"`
}

type CreatePhotoRequest struct {
	CategoryIds []uuid.UUID `json:"category_ids"`
}

type CreatePhotoResponse struct {
	Data PhotoData `json:"data"`
}

type UpdatePhotoRequest struct {
	IsUploaded *bool `json:"is_uploaded"`
}

type UpdatePhotoResponse struct {
	Data PhotoData `json:"data"`
}

type ListPhotoParams struct {
	Category *string `form:"category,omitempty" binding:"required"`
	// if specified, sort parameter should be either "star" or "created_at"
	// defaults to "created_at" if unspecified
	Sort      *string `form:"sort,omitempty"`
	Thumbnail *string `form:"thumbnail,omitempty"`
	PaginationParams
}

type ListMyPhotoParams struct {
	PaginationParams
}

// Returns an entries link to reveal other possible state transitions.
type ListPhotoLinks struct {
	Next string `json:"next"`
	//TODO entries
}

type ListPhotoResponse struct {
	Data  []PhotoListItemData `json:"data"`
	Meta  PhotoListMeta       `json:"meta"`
	Links ListPhotoLinks      `json:"links"`
}

type ListMyPhotoResponse struct {
	Data  []PhotoListItemData `json:"data"`
	Meta  MyPhotoListMeta     `json:"meta"`
	Links ListPhotoLinks      `json:"links"`
}

type GetPhotoRequest struct {
}

type GetPhotoResponse struct {
	Data PhotoData `json:"data"`
}

type DeletePhotoRequest struct {
}

type DeletePhotoResponse struct {
	Data PhotoData `json:"data"`
}
