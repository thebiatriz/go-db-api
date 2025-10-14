package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
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
