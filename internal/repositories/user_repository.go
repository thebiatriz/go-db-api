package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/thebiatriz/go-db-api/internal/models"
)

var ErrEmailAlreadyExists = errors.New("o email inserido já está cadastrado")
var ErrUserNotFound = errors.New("o usuário não foi encontrado na base de dados")

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{
		connection: connection,
	}
}

func (ur *UserRepository) GetUsers() ([]models.User, error) {
	query := "SELECT id, name, email FROM users"
	rows, err := ur.connection.Query(query)

	if err != nil {
		fmt.Println(err)
		return []models.User{}, err
	}

	var userList []models.User
	var userObj models.User

	for rows.Next() {
		err = rows.Scan(
			&userObj.ID,
			&userObj.Name,
			&userObj.Email,
		)

		if err != nil {
			fmt.Println(err)
			return []models.User{}, err
		}

		userList = append(userList, userObj)
	}

	rows.Close()

	return userList, nil
}

func (ur UserRepository) GetUserById(id_user int) (*models.User, error) {
	var user models.User

	query := "SELECT id, name, email FROM users WHERE id = $1"

	err := ur.connection.QueryRow(query, id_user).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

func (ur UserRepository) CreateUser(user models.User) (int, error) {
	var id int

	query := "INSERT INTO users(name, email, password) VALUES ($1, $2, $3) RETURNING id"

	err := ur.connection.QueryRow(query, user.Name, user.Email, user.Password).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, ErrEmailAlreadyExists
			}
		}
		return 0, nil
	}

	return id, nil
}

func (ur UserRepository) UpdateUser(user models.User) error {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"

	result, err := ur.connection.Exec(query, user.Name, user.Email, user.ID)

	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
