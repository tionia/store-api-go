package services

import (
	"store-api-go/internal/models"
	"store-api-go/internal/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepo
}

func NewTransactionService(repo *repositories.TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) Report(startDate string, endDate string) (*models.ReportResponse, error) {
	return s.repo.Report(startDate, endDate)
}
