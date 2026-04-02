package models

type CreateUserRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID string
	Name string
	Email string
}

type User struct {
	ID string
	Name string
	Email string
}