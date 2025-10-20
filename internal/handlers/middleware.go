package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thebiatriz/go-db-api/internal/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := []byte(os.Getenv("SECRET_JWT"))
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response := models.Response{
				Message: "O cabeçalho de autorização está vazio",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			response := models.Response{
				Message: "O cabeçalho de autorização está com a formatação inválida",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			response := models.Response{
				Message: "O token informado não é válido",
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Message = "Token expirado"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := int(claims["user_id"].(float64))
			userRole := claims["role"].(string)

			c.Set("userId", userId)
			c.Set("userRole", userRole)

			c.Next()

		} else {
			response := models.Response{
				Message: "Token inválido ou claims não encontrados",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
	}
}
