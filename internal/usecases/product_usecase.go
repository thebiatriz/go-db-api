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

func (pu *ProductUsecase) CreateProduct(product models.Product) (models.Product, error) {
	productId, err := pu.productRepository.CreateProduct(product)

	if err != nil {
		return models.Product{}, err
	}

	product.ID = productId

	return product, nil
}

func (pu *ProductUsecase) GetProductById(id_product int) (*models.Product, error) {
	product, err := pu.productRepository.GetProductById(id_product)

	if err != nil {
		return nil, err
	}

	return product, nil
}