package models

type User struct {
	UserID    uint64 `json:"userID"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type FormUser struct {
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type UserResponse struct {
	User User `json:"user"`
}
