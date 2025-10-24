package models

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required,min=2"`
	Price float64 `json:"price" binding:"required,gt=0"`
}