package handlers

import (
	"encoding/json"
	"net/http"
	"store-api-go/internal/models"
	"store-api-go/internal/services"
	"time"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// /api/checkout
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Method not allowed",
		})
	}
}

// /api/report/hari-ini
func (h *TransactionHandler) HandleReportToday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.ReportToday(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Method not allowed",
		})
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var request models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid request body",
		})
		return
	}

	transaction, err := h.service.Checkout(request.Items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Checkout success",
		Data:    transaction,
	})

}

func (h *TransactionHandler) ReportToday(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	date := now.Format("2006-01-02")
	reports, err := h.service.Report(date, date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(reports)

}
