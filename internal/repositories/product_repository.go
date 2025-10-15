package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/thebiatriz/go-db-api/internal/models"
)

var ErrProductNotFound = errors.New("o produto não foi encontrado na base de dados")

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]models.Product, error) {
	query := "SELECT id, product_name, price FROM product"
	rows, err := pr.connection.Query(query)

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
			&productObj.Price)

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
	query, err := pr.connection.Prepare("INSERT INTO product" +
		"(product_name, price)" +
		"VALUES ($1, $2) RETURNING id")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()

	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*models.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM product WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var product models.Product

	err = query.QueryRow(id_product).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
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

func (pr *ProductRepository) DeleteProduct(id_product int) error {
	query, err := pr.connection.Prepare("DELETE FROM product WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return err
	}

	result, err := query.Exec(id_product)

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
