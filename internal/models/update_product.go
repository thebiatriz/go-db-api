package models

type UpdateProductRequest struct {
	Name  *string  `json:"name" binding:"omitempty,min=2"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
}
