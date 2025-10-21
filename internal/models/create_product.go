package models

type CreateProductRequest struct {
	Name string `json:"name" bindind:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}