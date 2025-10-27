package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	
	"github.com/thebiatriz/go-db-api/internal/models"
)

var ErrProductNotFound = errors.New("o produto não foi encontrado na base de dados ou não pertence ao usuário")

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts(limit, offset int) ([]models.Product, error) {
	query := "SELECT id, name, price, user_id FROM products LIMIT $1 OFFSET $2"
	rows, err := pr.connection.Query(query, limit, offset)

	if err != nil {
		fmt.Println(err)
		return []models.Product{}, err
	}

	var productList []models.Product
	var productObj models.Product

	for rows.Next() {
		err = rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price,
			&productObj.UserID)

		if err != nil {
			fmt.Println(err)
			return []models.Product{}, err
		}

		productList = append(productList, productObj)
	}

	rows.Close()

	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product models.Product) (int, error) {
	var id int
	query, err := pr.connection.Prepare("INSERT INTO products" +
		"(name, price, user_id)" +
		"VALUES ($1, $2, $3) RETURNING id")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price, product.UserID).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()

	return id, nil
}

func (pr *ProductRepository) GetProductById(productId int) (*models.Product, error) {
	query, err := pr.connection.Prepare("SELECT id, name, price, user_id FROM products WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var product models.Product

	err = query.QueryRow(productId).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	query.Close()

	return &product, nil
}

func (pr *ProductRepository) DeleteProduct(productId int) error {
	query, err := pr.connection.Prepare("DELETE FROM products WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return err
	}

	result, err := query.Exec(productId)

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
		return ErrProductNotFound
	}

	query.Close()

	return nil
}

func (pr *ProductRepository) UpdateProduct(product models.Product, requesterId int) error {
	query, err := pr.connection.Prepare("UPDATE products SET name = $1, price = $2 WHERE id = $3 AND user_id = $4")

	if err != nil {
		fmt.Println(err)
		return err
	}

	result, err := query.Exec(product.Name, product.Price, product.ID, requesterId)

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
		return ErrProductNotFound
	}

	query.Close()

	return nil
}
