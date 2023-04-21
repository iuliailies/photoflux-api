package categories

import (
	public "github.com/iuliailies/photo-flux/pck/photoflux"
)

func ListCategoriesError(status int, detail string) public.Error {
	return public.Error{
		Title:  "List Categories Failed",
		Status: status,
		Detail: detail,
	}
}
