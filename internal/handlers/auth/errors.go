package auth

import (
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func RegisterError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Register Failed",
		Status: status,
		Detail: detail,
	}
}

func LoginError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Login Failed",
		Status: status,
		Detail: detail,
	}
}

func RefreshError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Refresh Token Failed",
		Status: status,
		Detail: detail,
	}
}

func LogoutError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Logout Failed",
		Status: status,
		Detail: detail,
	}
}
