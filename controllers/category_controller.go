package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"projet-forum/models"
)

type CategoryController struct {
	DB *sql.DB
}

// ListCategories récupère la liste des catégories
func (c *CategoryController) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.ListCategories(c.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetCategory récupère une catégorie par son ID
func (c *CategoryController) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := models.GetCategory(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// CreateCategory crée une nouvelle catégorie
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	createdCategory, err := models.CreateCategory(c.DB, category.Name, category.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

// UpdateCategory met à jour une catégorie
func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.ID = id
	if err := category.UpdateCategory(c.DB); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory supprime une catégorie
func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteCategory(c.DB, id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
