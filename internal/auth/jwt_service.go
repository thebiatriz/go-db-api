package auth

import (
	"errors"
	"fmt"
	"os"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userId int, userRole string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService() (JWTService, error) {
	secret := os.Getenv("SECRET_JWT")

	if secret == "" {
		return nil, errors.New("chave secreta JWT não está definida no ambiente")
	}
	return &jwtService{
		secretKey: []byte(secret),
	}, nil
}

func (s *jwtService) GenerateToken(userId int, userRole string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	jwt.MapClaims{
		"user_id": userId,
		"role": userRole,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	tokenString, err := token.SignedString(s.secretKey)

	if err != nil {
		return "", fmt.Errorf("erro ao assinar o token %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura não compatível: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
}

