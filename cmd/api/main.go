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

	router.POST("/users", UserHandler.CreateUser)
	router.POST("/login", UserHandler.Login)

	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware())

	productRoutes := protected.Group("/products")
	productRoutes.GET("/", ProductHandler.GetProducts)
	productRoutes.GET("/:id", ProductHandler.GetProductById)
	productRoutes.POST("/", ProductHandler.CreateProduct)
	productRoutes.DELETE("/:id", ProductHandler.DeleteProduct)
	productRoutes.PUT("/:id", ProductHandler.UpdateProduct)

	userRoutes := protected.Group("/users")
	userRoutes.GET("/", UserHandler.GetUsers)
	userRoutes.GET("/:id", UserHandler.GetUserById)
	userRoutes.DELETE("/:id", UserHandler.DeleteUser)
	userRoutes.PUT("/:id", UserHandler.UpdateUser)

	router.Run(":8080")
}
