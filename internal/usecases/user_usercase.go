package usecases

import (
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
)

type UserUsecase struct {
	userRepository repositories.UserRepository
}

func NewUserUsecase(userRepository repositories.UserRepository) UserUsecase {
	return UserUsecase {
		userRepository: userRepository,
	}
}

func (uu *UserUsecase) GetUsers() ([]models.User, error) {
	return uu.userRepository.GetUsers()
}