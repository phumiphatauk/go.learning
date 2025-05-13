package user

type Pagination struct {
	Page          int    `query:"page"`
	Limit         int    `query:"limit"`
	Sort          string `query:"sort"`
	SortDirection string `query:"sort_direction"`
}

type GetUserList struct {
	Pagination
	FirstName *string `query:"first_name"`
}

type PaginationResponse struct {
	Total int64 `json:"total_count"`
}

type GetUserListResponse struct {
	PaginationResponse
	Data []User `json:"data"`
}

type CreateUser struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type UpdateUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Active    bool   `json:"active"`
}

type User struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
