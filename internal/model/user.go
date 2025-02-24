package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=8"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone"`
}

type CreateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserLoginResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserServiceResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email"`
}
