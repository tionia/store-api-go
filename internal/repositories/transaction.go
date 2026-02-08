package repositories

import (
	"database/sql"
	"fmt"
	"store-api-go/internal/models"
	"strings"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (repo *TransactionRepo) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	dbTransaction, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer dbTransaction.Rollback()

	// Build WHERE IN query to fetch all products at once
	paramPlaceholder := make([]string, len(items))
	args := make([]interface{}, len(items))
	for i, item := range items {
		paramPlaceholder[i] = fmt.Sprintf("$%d", i+1) // results in $1, $2, ...
		args[i] = item.ProductID
	}

	query := fmt.Sprintf("SELECT id, name, price, stock FROM products WHERE id IN (%s)", strings.Join(paramPlaceholder, ", "))
	rows, err := dbTransaction.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map products by ID for quick lookup
	productMap := make(map[int]models.Product)
	for rows.Next() {
		var productResult models.Product
		if err := rows.Scan(&productResult.ID, &productResult.Name, &productResult.Price, &productResult.Stock); err != nil {
			return nil, err
		}

		productMap[productResult.ID] = productResult
	}

	// Validate all products exist and build details
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("product id %d is out of stock", item.ProductID)
		}

		subtotal := product.Price * item.Quantity
		totalAmount += subtotal

		_, err = dbTransaction.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = dbTransaction.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)

	if err != nil {
		return nil, err
	}

	// Set TransactionID and build batch insert
	insertParamPlaceHolder := make([]string, len(details))
	insertArgs := make([]interface{}, 0, len(details)*4)

	for i := range details {
		details[i].TransactionID = transactionID
		insertParamPlaceHolder[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		insertArgs = append(insertArgs, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
	}

	// Batch insert query
	_, err = dbTransaction.Exec(
		fmt.Sprintf("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES %s", strings.Join(insertParamPlaceHolder, ", ")),
		insertArgs...,
	)
	if err != nil {
		return nil, err
	}

	// Loop the details and insert to db
	// for _, detail := range details {
	// 	detail.TransactionID = transactionID
	// 	_, err = dbTransaction.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
	// 		transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	if err := dbTransaction.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepo) Report(startDate string, endDate string) (*models.ReportResponse, error) {
	return &models.ReportResponse{}, nil
}
