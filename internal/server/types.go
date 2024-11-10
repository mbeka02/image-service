package server

type CreateUserRequest struct {
	Username string `json:"user_name" validate:"required"`
	Fullname string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	Username string `json:"user_name" validate:"required"`
	Fullname string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}
