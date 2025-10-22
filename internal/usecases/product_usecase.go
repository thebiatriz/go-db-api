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

func (pu *ProductUsecase) CreateProduct(product models.Product) (*models.Product, error) {
	productId, err := pu.productRepository.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	product.ID = productId

	return &product, nil
}

func (pu *ProductUsecase) GetProductById(id_product int) (*models.Product, error) {
	product, err := pu.productRepository.GetProductById(id_product)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pu *ProductUsecase) DeleteProduct(id_product int, requesterId int, requesterRole string) error {
	productToDelete, err := pu.productRepository.GetProductById(id_product)
	if err != nil {
		return err
	}
	if productToDelete == nil {
		return repositories.ErrProductNotFound
	}

	isAdmin := requesterRole == "admin"
	isOwner := requesterId == productToDelete.UserID

	if !isAdmin && !isOwner {
		return ErrNotAuthorized 
	}

	err = pu.productRepository.DeleteProduct(id_product)
	if err != nil {
		return err
	}

	return nil
}

func (pu *ProductUsecase) UpdateProduct(product models.Product, requesterId int) (*models.Product, error) {
	productToUpdate, err := pu.productRepository.GetProductById(product.ID)
	if err != nil {{
		return nil, err
	}}

	if productToUpdate == nil {
		return nil, repositories.ErrProductNotFound
	}

	isOwner := productToUpdate.UserID == requesterId

	if !isOwner {
		return nil, ErrNotAuthorized
	}

	err = pu.productRepository.UpdateProduct(product, requesterId)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
