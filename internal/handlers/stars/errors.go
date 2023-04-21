package stars

import (
	public "github.com/iuliailies/photo-flux/pck/photoflux"
)

func StarPhotoError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Star Photo Failed",
		Status: status,
		Detail: detail,
	}
}

func IsPhotoStarredError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Star Photo Failed",
		Status: status,
		Detail: detail,
	}
}
