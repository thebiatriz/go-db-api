package server

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/auth"
	"github.com/thebiatriz/go-db-api/internal/database"
	"github.com/thebiatriz/go-db-api/internal/handlers"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
)

var (
	errConnectDatabase = errors.New("erro ao conectar ao banco de dados")
	errInitializingJWT = errors.New("erro ao inicializar o serviço do JWT")
)

type Server struct {
	engine         *gin.Engine
	productHandler handlers.ProductHandler
	userHandler    handlers.UserHandler
	jwtService     auth.JWTService
	userRepository repositories.UserRepository
}

func NewServer() (*Server, error) {
	dbConnection, err := database.ConnectDB()

	if err != nil {
		return nil, errConnectDatabase
	}

	jwtService, err := auth.NewJWTService()

	if err != nil {
		return nil, errInitializingJWT
	}

	productRepository := repositories.NewProductRepository(dbConnection)
	productUseCase := usecases.NewProductUsecase(productRepository)
	productHandler := handlers.NewProductHandler(productUseCase)

	userRepository := repositories.NewUserRepository(dbConnection)
	userUseCase := usecases.NewUserUsecase(userRepository, jwtService)
	userHandler := handlers.NewUserHandler(userUseCase)

	router := gin.Default()

	server := &Server{
		engine:         router,
		productHandler: productHandler,
		userHandler:    userHandler,
		jwtService:     jwtService,
		userRepository: userRepository,
	}

	server.registerRoutes()
	return server, nil
}

func (s *Server) registerRoutes() {
	s.engine.POST("/users", s.userHandler.CreateUser)
	s.engine.POST("/login", s.userHandler.Login)

	protected := s.engine.Group("/")
	protected.Use(handlers.AuthMiddleware(s.jwtService, s.userRepository))
	
	productRoutes := protected.Group("/products")
	productRoutes.GET("/", s.productHandler.GetProducts)
	productRoutes.GET("/:id", s.productHandler.GetProductById)
	productRoutes.POST("/", s.productHandler.CreateProduct)
	productRoutes.DELETE("/:id", s.productHandler.DeleteProduct)
	productRoutes.PUT("/:id", s.productHandler.UpdateProduct)

	userRoutes := protected.Group("/users")
	userRoutes.GET("/me", s.userHandler.GetMe)
	userRoutes.GET("/", s.userHandler.GetUsers)
	userRoutes.GET("/:id", s.userHandler.GetUserById)
	userRoutes.DELETE("/me", s.userHandler.DeleteMe)
	userRoutes.DELETE("/:id", s.userHandler.DeleteUser)
	userRoutes.PUT("/me", s.userHandler.UpdateMe)
	userRoutes.PUT("/:id", s.userHandler.UpdateUser)
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
