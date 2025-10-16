package handlers

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
)

type UserHandler struct {
	userUsecase usecases.UserUsecase
}

func NewUserHandler(userUsecase usecases.UserUsecase) UserHandler {
	return UserHandler{
		userUsecase: userUsecase,
	}
}

func (u UserHandler) GetUsers(c *gin.Context) {
	users, err := u.userUsecase.GetUsers()

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (u UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := u.userUsecase.CreateUser(user)

	if err != nil {
		if errors.Is(err, repositories.ErrEmailAlreadyExists) {
			response := models.Response{
				Message: "O email inserido já está cadastrado na base de dados",
			}
			c.IndentedJSON(http.StatusConflict, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}
