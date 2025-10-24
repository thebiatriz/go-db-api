package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"github.com/thebiatriz/go-db-api/internal/usecases"
)

const (
	msgBadRequest          = "Ocorreu um erro ao receber os dados na requisição"
	msgInternalServerError = "Ocorreu um erro interno no servidor"
	msgUnauthorized        = "Você não tem permissão para realizar essa ação"
	msgInvalidIDNumeric    = "O id precisa ser um número"
	msgInvalidIDFormat     = "Formato inválido do id do usuário"
	msgUserIDMissing       = "Não foi possível identificar o usuário"
	msgInvalidRoleFormat   = "Formato inválido da permissão do usuário"
	msgUserRoleMissing     = "Não foi possível identificar a permissão do usuário"
)

func handlePayloadErrors(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	var wrongTag string

	if errors.As(err, &validationErrs) {
		firstError := validationErrs[0]

		switch firstError.Tag() {
		case "min":
			wrongTag = "tamanho mínimo insuficiente"
		case "email":
			wrongTag = "formato do email inválido"
		case "required":
			wrongTag = "campo obrigatório ausente"
		case "gt":
			wrongTag = "valor precisa ser maior que zero"
		default:
			wrongTag = fmt.Sprintf("regra '%s'", firstError.Tag())
		}

		errorMessage := fmt.Sprintf("Erro no campo '%s': falhou na regra '%s'",
			firstError.Field(),
			wrongTag,
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{Message: errorMessage})
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{Message: msgBadRequest})
}

func handleDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecases.ErrUserNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, models.Response{Message: "O usuário não foi encontrado na base de dados"})

	case errors.Is(err, usecases.ErrProductNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, models.Response{Message: "O produto não foi encontrado na base de dados"})

	case errors.Is(err, usecases.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Message: "Senha inválida"})

	case errors.Is(err, usecases.ErrNotAuthorized):
		c.AbortWithStatusJSON(http.StatusForbidden, models.Response{Message: msgUnauthorized})

	case errors.Is(err, repositories.ErrEmailAlreadyExists):
		c.AbortWithStatusJSON(http.StatusConflict, models.Response{Message: "O email inserido já está cadastrado na base de dados"})

	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Message: msgInternalServerError})
	}
}
