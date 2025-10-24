package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
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
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

func (p *productHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	requesterIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}

		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterId, ok := requesterIdStr.(int)
	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	err := c.BindJSON(&req)

	if err != nil {
		handlePayloadErrors(c, err)
		return
	}

	productToCreate := models.Product{
		Name:   req.Name,
		Price:  req.Price,
		UserID: requesterId,
	}

	insertedProduct, err := p.productUsecase.CreateProduct(productToCreate)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, insertedProduct)
}

func (p *productHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: msgInvalidIDNumeric,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.GetProductById(productId)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func (p *productHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: msgInvalidIDNumeric,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	requesterIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}
	
	requesterId, ok := requesterIdStr.(int)
	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRoleStr, exists := c.Get("userRole")
	if !exists {
		response := models.Response{
			Message: msgUserRoleMissing,
		}

		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRole, ok := requesterRoleStr.(string)
	if !ok {
		response := models.Response{
			Message: msgInvalidRoleFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	err = p.productUsecase.DeleteProduct(productId, requesterId, requesterRole)

	if err != nil {
		handleDomainError(c, err)
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
			Message: msgInvalidIDNumeric,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = c.BindJSON(&req)

	if err != nil {
		handlePayloadErrors(c, err)
		return
	}

	requesterIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterId, ok := requesterIdStr.(int)
	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	updatedProduct, err := p.productUsecase.UpdateProduct(targetId, requesterId, req)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedProduct)
}
