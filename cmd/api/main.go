package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/auth"
	"github.com/thebiatriz/go-db-api/internal/database"
	"github.com/thebiatriz/go-db-api/internal/handlers"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
	"log"
)

func main() {
	router := gin.Default()

	dbConnection, err := database.ConnectDB()

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	jwtService, err := auth.NewJWTService()

	if err != nil {
		log.Fatalf("Erro ao inicializar o serviço do JWT: %v", err)
	}

	ProductRepository := repositories.NewProductRepository(dbConnection)
	ProductUseCase := usecases.NewProductUsecase(ProductRepository)
	ProductHandler := handlers.NewProductHandler(ProductUseCase)

	UserRepository := repositories.NewUserRepository(dbConnection)
	UserUsecase := usecases.NewUserUsecase(UserRepository, jwtService)
	UserHandler := handlers.NewUserHandler(UserUsecase)

	router.POST("/users", UserHandler.CreateUser)
	router.POST("/login", UserHandler.Login)

	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware(jwtService, UserRepository))

	productRoutes := protected.Group("/products")
	productRoutes.GET("/", ProductHandler.GetProducts)
	productRoutes.GET("/:id", ProductHandler.GetProductById)
	productRoutes.POST("/", ProductHandler.CreateProduct)
	productRoutes.DELETE("/:id", ProductHandler.DeleteProduct)
	productRoutes.PUT("/:id", ProductHandler.UpdateProduct)

	userRoutes := protected.Group("/users")
	userRoutes.GET("/me", UserHandler.GetMe)
	userRoutes.GET("/", UserHandler.GetUsers)
	userRoutes.GET("/:id", UserHandler.GetUserById)
	userRoutes.DELETE("/me", UserHandler.DeleteMe)
	userRoutes.DELETE("/:id", UserHandler.DeleteUser)
	userRoutes.PUT("/me", UserHandler.UpdateMe)
	userRoutes.PUT("/:id", UserHandler.UpdateUser)

	router.Run(":8080")
}