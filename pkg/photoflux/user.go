package photoflux

const UserType = "user"

// Data returned about a user when a single one is returned.
type UserData struct {
	ResourceID
	Attributes UserAttributes `json:"attributes"`
	Links      UserLinks      `json:"links"`
}

type UserAttributes struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Timestamps
}

type UserLinks struct {
	Self string `json:"self"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Data UserData `json:"data"`
}
