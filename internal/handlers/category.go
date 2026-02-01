package handlers

import (
	"encoding/json"
	"go-categories-api/internal/models"
	"go-categories-api/internal/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// handle /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.GetAll(w)
	case http.MethodPost:
		h.Create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Method not allowed",
		})
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handle /api/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get id from path param
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid ID",
		})
		// http.Error(w, "Invalid ID", http.StatusBadRequest)
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
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter) {
	categories, err := h.service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Categories retrieved",
		Data:    categories,
	})
	// json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid request body",
		})

		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&newCategory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})

		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Category created",
		Data:    newCategory,
	})

	// json.NewEncoder(w).Encode(newCategory)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, id int) {
	category, err := h.service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})

		// http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Category retrieved",
		Data:    category,
	})

	// json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var categoryUpdate models.Category
	err := json.NewDecoder(r.Body).Decode(&categoryUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: "Invalid request body",
		})

		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	categoryUpdate.ID = id
	err = h.service.Update(&categoryUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})

		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Category updated",
		Data:    categoryUpdate,
	})
	// json.NewEncoder(w).Encode(categoryUpdate)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, id int) {
	err := h.service.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  "FAIL",
			Message: err.Error(),
		})

		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "OK",
		Message: "Category deleted",
	})
}
