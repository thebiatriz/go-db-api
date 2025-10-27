package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/usecases"
)

type ProductHandler struct {
	productUsecase usecases.ProductUsecase
}

func NewProductHandler(usecase usecases.ProductUsecase) ProductHandler {
	return ProductHandler{
		productUsecase: usecase,
	}
}

func (p *ProductHandler) GetProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	queryName := c.DefaultQuery("name", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	products, err := p.productUsecase.GetProducts(page, limit, queryName)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

func (p *ProductHandler) CreateProduct(c *gin.Context) {
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

	req.Name = strings.TrimSpace(req.Name)

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

func (p *ProductHandler) GetProductById(c *gin.Context) {
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

func (p *ProductHandler) DeleteProduct(c *gin.Context) {
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

func (p *ProductHandler) UpdateProduct(c *gin.Context) {
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

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
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
