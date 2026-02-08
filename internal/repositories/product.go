package repositories

import (
	"database/sql"
	"errors"
	"store-api-go/internal/models"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (repo *ProductRepo) GetAll(name string) ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products"

	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repo *ProductRepo) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock).Scan(&product.ID)

	return err
}

func (repo *ProductRepo) GetByID(id int) (*models.Product, error) {
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"

	var product models.Product
	err := repo.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (repo *ProductRepo) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return nil
}

func (repo *ProductRepo) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return err
}
