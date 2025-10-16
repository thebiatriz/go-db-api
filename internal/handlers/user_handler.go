package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
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
	}

	c.IndentedJSON(http.StatusOK, users)
}