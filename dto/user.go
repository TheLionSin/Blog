package dto

import "Blog/models"

type RegisterInput struct {
	Nickname string `json:"nickname" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserInput struct {
	Nickname string `json:"nickname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty, min=5"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
}

func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
	}
}

type CreateUserInput struct {
	Nickname string `json:"nickname" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}
