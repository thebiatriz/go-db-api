package models

type UpdateUserRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=3"`
	Email    *string `json:"email" binding:"omitempty,email"`
}