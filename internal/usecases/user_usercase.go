package usecases

import (
	"errors"
	"fmt"
	"github.com/thebiatriz/go-db-api/internal/auth"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("o usuário não existe na base de dados")
	ErrInvalidCredentials = errors.New("senha inválida")
	ErrNotAuthorized      = errors.New("não autorizado para modificação")
)

type UserUsecase struct {
	userRepository repositories.UserRepository
	jwtService auth.JWTService
}

func NewUserUsecase(userRepository repositories.UserRepository, jwtService auth.JWTService) UserUsecase {
	return UserUsecase{
		userRepository: userRepository,
		jwtService: jwtService,
	}
}

func (uu *UserUsecase) Login(userEmail string, password string) (string, error) {
	user, err := uu.userRepository.GetUserByEmail(userEmail)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := uu.jwtService.GenerateToken(user.ID, user.Role)

	if err != nil {
		return "", fmt.Errorf("erro interno ao gerar o token: %w", err)
	}

	return token, nil
}

func (uu *UserUsecase) GetUsers() ([]models.User, error) {
	return uu.userRepository.GetUsers()
}

func (uu UserUsecase) GetUserById(userId int) (*models.User, error) {
	user, err := uu.userRepository.GetUserById(userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
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

func (uu UserUsecase) DeleteUser(targetUserId, requesterId int, requesterRole string) error {
	isAdmin := requesterRole == "admin"
	isSelf := targetUserId == requesterId

	if !isAdmin && !isSelf {
		return ErrNotAuthorized
	}

	err := uu.userRepository.DeleteUser(targetUserId)

	if err != nil {
		return err
	}

	return nil
}

func (uu UserUsecase) UpdateUser(targetUserId, requesterId int, requesterRole string, req models.UpdateUserRequest) (*models.User, error) {
	currentUser, err := uu.userRepository.GetUserById(targetUserId)
	if err != nil {
		return nil, err
	}

	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	isSelf := requesterId == currentUser.ID
	isAdmin := requesterRole == "admin"

	if !isAdmin && !isSelf {
		return nil, ErrNotAuthorized
	}

	if req.Name != nil {
		currentUser.Name = *req.Name
	}

	if req.Email != nil {
		currentUser.Email = *req.Email
	}

	err = uu.userRepository.UpdateUser(*currentUser)

	if err != nil {
		return nil, err
	}

	return currentUser, nil
}
