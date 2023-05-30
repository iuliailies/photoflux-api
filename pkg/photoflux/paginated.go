package photoflux

type PaginationParams struct {
	After *string `form:"after,omitempty"`
	Limit *int    `form:"limit,omitempty"`
}
