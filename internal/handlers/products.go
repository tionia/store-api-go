package handlers

import (
	"encoding/json"
	"net/http"
	"store-api-go/internal/models"
	"store-api-go/internal/services"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// handle /api/products
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Method not allowed",
		})
	}
}

// handle /api/categories/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get id from path param
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid ID",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Method not allowed",
		})
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.service.GetAll(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Products retrieved",
		Data:    products,
	})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid request body",
		})

		return
	}

	err = h.service.Create(&newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Product created",
		Data:    newProduct,
	})

}

func (h *ProductHandler) GetByID(w http.ResponseWriter, id int) {
	product, err := h.service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Product retrieved",
		Data:    product,
	})

}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var productUpdate models.Product
	err := json.NewDecoder(r.Body).Decode(&productUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid request body",
		})
		return
	}

	productUpdate.ID = id
	err = h.service.Update(&productUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Product updated",
		Data:    productUpdate,
	})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, id int) {
	err := h.service.Delete(id)
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
		Message: "Product deleted",
	})
}
