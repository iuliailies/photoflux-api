package boards

import (
	public "github.com/iuliailies/photo-flux/pkg/photoflux"
)

func CreateBoardError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create Board Failed",
		Status: status,
		Detail: detail,
	}
}

func UpdateBoardError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Update Board Failed",
		Status: status,
		Detail: detail,
	}
}

func ListBoardError(status int, detail string) public.Error {
	return public.Error{
		Title:  "List Board Failed",
		Status: status,
		Detail: detail,
	}
}

func GetBoardError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Get Board Failed",
		Status: status,
		Detail: detail,
	}
}
