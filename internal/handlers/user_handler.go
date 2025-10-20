package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func (u UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	err := c.BindJSON(&req)

	if err != nil {
		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	token, err := u.userUsecase.Login(req.Email, req.Password)

	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			response := models.Response{
				Message: "Senha inválida",
			}
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		if errors.Is(err, usecases.ErrUserNotFound) {
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

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
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
	var req models.CreateUserRequest

	err := c.BindJSON(&req)

	if err != nil {
		var validationErrs validator.ValidationErrors
		var wrongTag string

		if errors.As(err, &validationErrs) {
			firstError := validationErrs[0]

			switch firstError.Tag() {
			case "min":
				wrongTag = "tamanho mínimo"
			case "email":
				wrongTag = "formato do email"
			case "required": 
				wrongTag = "campo obrigatório"
			default: 
				wrongTag = fmt.Sprintf("regra '%s'", firstError.Tag())
			}

			errorMessage := fmt.Sprintf("Erro no campo '%s': falhou na regra '%s'", firstError.Field(), wrongTag)

			response := models.Response{
				Message: errorMessage,
			}
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		response := models.Response{
			Message: "Ocorreu um erro ao receber os dados na requisição",
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	userToCreate := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	newUser, err := u.userUsecase.CreateUser(userToCreate)

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
