package main

import (
	"encoding/json"
	"net/http"
	"store-api-go/internal/models"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "Server is running",
	})
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category

	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func getCategories(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func getCategoryById(w http.ResponseWriter, id int) {
	for _, category := range categories {
		if category.ID == id {
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var categoryUpdate models.Category

	err := json.NewDecoder(r.Body).Decode(&categoryUpdate)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i, category := range categories {
		if category.ID == id {
			categoryUpdate.ID = id
			categories[i] = categoryUpdate // Modify the original slice using index
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, id int) {
	for i, category := range categories {
		if category.ID == id {
			var remainingCategories []models.Category

			remainingCategories = append(categories[:i], categories[i+1:]...)

			categories = remainingCategories

			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "Category deleted",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}
