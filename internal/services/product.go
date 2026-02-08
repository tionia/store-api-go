package services

import (
	"store-api-go/internal/models"
	"store-api-go/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepo
}

func NewProductService(repo *repositories.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(Product *models.Product) error {
	return s.repo.Update(Product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
