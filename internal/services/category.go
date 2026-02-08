package services

import (
	"store-api-go/internal/models"
	"store-api-go/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepo
}

func NewCategoryService(repo *repositories.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(name string) ([]models.Category, error) {
	return s.repo.GetAll(name)
}

func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
