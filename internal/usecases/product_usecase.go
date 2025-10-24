package usecases

import (
	"errors"

	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
)

var ErrProductNotFound = errors.New("o usuário não existe na base de dados")

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

func (pu *ProductUsecase) GetProductById(productId int) (*models.Product, error) {
	product, err := pu.productRepository.GetProductById(productId)

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	return product, nil
}

func (pu *ProductUsecase) DeleteProduct(productId int, requesterId int, requesterRole string) error {
	productToDelete, err := pu.productRepository.GetProductById(productId)
	if err != nil {
		return err
	}
	if productToDelete == nil {
		return ErrProductNotFound
	}

	isAdmin := requesterRole == "admin"
	isOwner := requesterId == productToDelete.UserID

	if !isAdmin && !isOwner {
		return ErrNotAuthorized
	}

	err = pu.productRepository.DeleteProduct(productId)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}

func (pu *ProductUsecase) UpdateProduct(targetId, requesterId int, req models.UpdateProductRequest) (*models.Product, error) {
	productToUpdate, err := pu.productRepository.GetProductById(targetId)
	if err != nil {
		return nil, err
	}

	if productToUpdate == nil {
		return nil, ErrProductNotFound
	}

	isOwner := productToUpdate.UserID == requesterId

	if !isOwner {
		return nil, ErrNotAuthorized
	}

	if req.Name != nil {
		productToUpdate.Name = *req.Name
	}

	if req.Price != nil {
		productToUpdate.Price = *req.Price
	}

	err = pu.productRepository.UpdateProduct(*productToUpdate, requesterId)

	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return productToUpdate, nil
}
