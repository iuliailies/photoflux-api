package users

import (
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func CreateUserError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create User Failed",
		Status: status,
		Detail: detail,
	}
}
