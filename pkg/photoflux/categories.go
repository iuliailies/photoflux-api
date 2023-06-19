package photoflux

import "strings"

const CategoryType = "category"

type CategoryData struct {
	ResourceID
	Attributes CategoryAttributes `json:"attributes"`
	Links      CategoryLinks      `json:"links"`
}

// Returns links to reveal other possible state transitions.
type ListCategoryLinks struct {
	Self string `json:"self"`
	//TODO entries
}

type ListCategoryResponse struct {
	Data  []CategoryListItemData `json:"data"`
	Links ListCategoryLinks      `json:"links"`
}

type CategoryListItemData struct {
	ResourceID
	Attributes CategoryAttributes    `json:"attributes"`
	Links      CategoryListItemLinks `json:"links"`
}

type CategoryAttributes struct {
	Name string `json:"name"`
	Timestamps
}

type CategoryLinks struct {
	Self string `json:"self"`
}

type CategoryListItemLinks struct {
	Self string `json:"self"`
}

func CategoriesFromURL(urlSnippet string) []string {
	return strings.Split(urlSnippet, ".")
}
