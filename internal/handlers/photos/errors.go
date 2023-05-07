package photos

import (
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func CreatePhotoError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create Photo Failed",
		Status: status,
		Detail: detail,
	}
}

func DeletePhotoError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Delete Photo Failed",
		Status: status,
		Detail: detail,
	}
}

func ListPhotoError(status int, detail string) public.Error {
	return public.Error{
		Title:  "List Photo Failed",
		Status: status,
		Detail: detail,
	}
}

func UpdatePhotoError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Update Photo Failed",
		Status: status,
		Detail: detail,
	}
}

func GetPhotoError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Get Photo Failed",
		Status: status,
		Detail: detail,
	}
}
