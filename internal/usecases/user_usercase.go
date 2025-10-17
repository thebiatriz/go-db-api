package usecases

import (
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository repositories.UserRepository
}

func NewUserUsecase(userRepository repositories.UserRepository) UserUsecase {
	return UserUsecase{
		userRepository: userRepository,
	}
}

func (uu *UserUsecase) GetUsers() ([]models.User, error) {
	return uu.userRepository.GetUsers()
}

func (uu UserUsecase) GetUserById(id_user int) (*models.User, error) {
	user, err := uu.userRepository.GetUserById(id_user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *UserUsecase) CreateUser(user models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	userId, err := uu.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	user.ID = userId

	return &user, nil
}

func (uu UserUsecase) DeleteUser(id_user int) error {
	err := uu.userRepository.DeleteUser(id_user)

	if err != nil {
		return err
	}

	return nil
}

func (uu UserUsecase) UpdateUser(user models.User) (*models.User, error) {
	err := uu.userRepository.UpdateUser(user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}