package photoflux

const BoardType = "board"

type BoardData struct {
	ResourceID
	Timestamps
	UserId string `json:"user_id"`
	Data   string `json:"data"`
}

type CreateBoardRequest struct {
	Data string `json:"data"`
}

type CreateBoardResponse struct {
	Data BoardData `json:"data"`
}

type UpdateBoardRequest struct {
	Data string `json:"data"`
}

type UpdateBoardResponse struct {
	Data string `json:"data"`
}

type GetBoardRequest struct {
}

type GetBoardResponse struct {
	Data BoardData `json:"data"`
}

type ListBoardsRequest struct {
}

type ListBoardsResponse struct {
	Data []BoardData `json:"data"`
}
