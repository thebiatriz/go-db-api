package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
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
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusOK, products)
}

func (p *productHandler) CreateProduct(c *gin.Context) {
	var product models.Product

	err := c.BindJSON(&product)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	insertedProduct, err := p.productUsecase.CreateProduct(product)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
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

	err = p.productUsecase.DeleteProduct(productId)

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
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}

	c.Status(http.StatusNoContent)
}
