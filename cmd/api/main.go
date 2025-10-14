package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/database"
	"github.com/thebiatriz/go-db-api/internal/handlers"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
)

func main() {
	router := gin.Default()

	dbConnection, err := database.ConnectDB()

	if err != nil {
		panic(err)
	}

	//Camada de repository
	ProductRepository := repositories.NewProductRepository(dbConnection)

	//Camada de usecase
	ProductUseCase := usecases.NewProductUsecase(ProductRepository)

	//Camada de handlers
	ProductHandler := handlers.NewProductHandler(ProductUseCase)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/products", ProductHandler.GetProducts)

	router.Run(":8080")
}
