package repositories

import (
	"database/sql"
	"errors"
	"store-api-go/internal/models"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (repo *CategoryRepo) GetAll(name string) ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories"

	// data type can be any type --> use interface. but the interface can be multiple??
	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...) // adding argument to query
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (repo *CategoryRepo) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)

	return err
}

func (repo *CategoryRepo) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var category models.Category
	err := repo.db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("Category not found")
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepo) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Category not found")
	}

	return nil
}

func (repo *CategoryRepo) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Category not found")
	}

	return err
}
