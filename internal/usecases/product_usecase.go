package usecases

import (
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
)

type ProductUsecase struct {
	productRepository repositories.ProductRepository
}

func NewProductUsecase(repo repositories.ProductRepository) ProductUsecase {
	return ProductUsecase{
		productRepository: repo,
	}
}

func (pu *ProductUsecase) GetProducts() ([]models.Product, error) {
	return pu.productRepository.GetProducts()
}