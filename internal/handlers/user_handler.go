package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
	"net/http"
	"strconv"
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

func (u UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response := models.Response{
			Message: "Id não pode estar vazio",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	userId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id precisa ser um número",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	user, err := u.userUsecase.GetUserById(userId)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	if user == nil {
		response := models.Response{
			Message: "O usuário não foi encontrado na base de dados",
		}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
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

func (u UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response := models.Response{
			Message: "Id não pode estar vazio",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	userId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id precisa ser um número",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = u.userUsecase.DeleteUser(userId)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			response := models.Response{
				Message: "O usuário não foi encontrado na base de dados",
			}
			c.IndentedJSON(http.StatusNotFound, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.Status(http.StatusNoContent)

}

func (u UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if id == "" {
		response := models.Response{
			Message: "Id não pode estar vazio",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	userId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: "Id precisa ser um número",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = c.BindJSON(&user)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	user.ID = userId
	updatedUser, err := u.userUsecase.UpdateUser(user)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			response := models.Response{
				Message: "O usuário não foi encontrado na base de dados",
			}
			c.IndentedJSON(http.StatusNotFound, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro interno no servidor",
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}
