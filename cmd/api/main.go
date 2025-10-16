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

	ProductRepository := repositories.NewProductRepository(dbConnection)
	ProductUseCase := usecases.NewProductUsecase(ProductRepository)
	ProductHandler := handlers.NewProductHandler(ProductUseCase)

	UserRepository := repositories.NewUserRepository(dbConnection)
	UserUsecase := usecases.NewUserUsecase(UserRepository)
	UserHandler := handlers.NewUserHandler(UserUsecase)

	router.GET("/products", ProductHandler.GetProducts)
	router.GET("/products/:id", ProductHandler.GetProductById)
	router.POST("/products", ProductHandler.CreateProduct)
	router.DELETE("/products/:id", ProductHandler.DeleteProduct)
	router.PUT("/products/:id", ProductHandler.UpdateProduct)

	router.GET("/users", UserHandler.GetUsers)
	router.POST("/users", UserHandler.CreateUser)

	router.Run(":8080")
}
