package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
	"net/http"
	"strconv"
)

type productHandler struct {
	productUsecase usecases.ProductUsecase
}

func NewProductHandler(usecase usecases.ProductUsecase) productHandler {
	return productHandler{
		productUsecase: usecase,
	}
}

func (p *productHandler) GetProducts(c *gin.Context) {
	products, err := p.productUsecase.GetProducts()
	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

func (p *productHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	requestIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: "Não foi possível identificar o usuário",
		}

		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requestId := requestIdStr.(int)

	err := c.BindJSON(&req)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	productToCreate := models.Product{
		Name:   req.Name,
		Price:  req.Price,
		UserID: requestId,
	}

	insertedProduct, err := p.productUsecase.CreateProduct(productToCreate)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusCreated, insertedProduct)
}

func (p *productHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response := models.Response{
			Message: "Id do produto não pode ser nulo",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id do produto precisa ser um número",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.GetProductById(productId)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := models.Response{
			Message: "O produto não foi encontrado na base de dados",
		}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func (p *productHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response := models.Response{
			Message: "Id do produto não pode ser nulo",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	productId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id do produto precisa ser um número",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	requesterIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: "Não foi possível identificar o usuário",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}
	requesterId := requesterIdStr.(int)

	requesterRoleStr, exists := c.Get("userRole")
	if !exists {
		response := models.Response{
			Message: "Não foi possível identificar a permissão do usuário",
		}

		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRole := requesterRoleStr.(string)

	err = p.productUsecase.DeleteProduct(productId, requesterId, requesterRole)

	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			response := models.Response{
				Message: "O produto não foi encontrado na base de dados",
			}
			c.IndentedJSON(http.StatusNotFound, response)
			return
		}

		if errors.Is(err, usecases.ErrNotAuthorized) {
			response := models.Response{
				Message: "Você não tem permissão para deletar esse produto",
			}
			c.IndentedJSON(http.StatusUnauthorized, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}

	c.Status(http.StatusNoContent)
}

func (p *productHandler) UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if id == "" {
		response := models.Response{
			Message: "Id do produto não pode ser nulo",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id do produto precisa ser um número",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = c.BindJSON(&product)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	product.ID = productId

	updatedProduct, err := p.productUsecase.UpdateProduct(product)

	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			response := models.Response{
				Message: "O produto não foi encontrado na base de dados",
			}
			c.IndentedJSON(http.StatusNotFound, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}

		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedProduct)
}
