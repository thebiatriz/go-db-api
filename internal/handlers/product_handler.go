package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		var validationErrs validator.ValidationErrors
		var wrongTag string

		if errors.As(err, &validationErrs) {
			firstError := validationErrs[0]

			switch firstError.Tag() {
			case "min":
				wrongTag = "tamanho mínimo insuficiente"
			case "gt":
				wrongTag = "preço precisa ser maior que zero"
			case "required":
				wrongTag = "campo obrigatório ausente"
			default:
				wrongTag = fmt.Sprintf("regra '%s'", firstError.Tag())
			}

			errorMessage := fmt.Sprintf("Erro no campo '%s': falha na regra '%s'", firstError.Field(), wrongTag)

			response := models.Response{
				Message: errorMessage,
			}
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

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
		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
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
	var req models.UpdateProductRequest
	id := c.Param("id")

	targetId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id do produto precisa ser um número",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = c.BindJSON(&req)

	if err != nil {
		var validationErrs validator.ValidationErrors
		var wrongTag string

		if errors.As(err, &validationErrs) {
			firstError := validationErrs[0]

			switch firstError.Tag() {
			case "min":
				wrongTag = "tamanho mínimo insuficiente"
			case "gt":
				wrongTag = "preço precisa ser maior que zero"
			default:
				wrongTag = fmt.Sprintf("regra '%s'", firstError.Tag())
			}

			errorMessage := fmt.Sprintf("Erro no campo '%s': falha na regra '%s'", firstError.Field(), wrongTag)

			response := models.Response{
				Message: errorMessage,
			}
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
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

	updatedProduct, err := p.productUsecase.UpdateProduct(targetId, requesterId, req)

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
				Message: "Você não tem permissão para atualizar esse produto",
			}
			c.IndentedJSON(http.StatusUnauthorized, response)
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
