package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// HANDLERS DE CATEGORY

// POST /categories: Crea una categoría nueva
func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "El nombre de la categoría no puede estar vacío", http.StatusBadRequest)
		return
	}

	category, err := s.queries.CreateCategory(r.Context(), input.Name)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GET /categories: Trae todas las categorías
func (s *Server) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := s.queries.ListCategories(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las categorías", http.StatusInternalServerError)
		return
	}

	if categories == nil {
		categories = []sqlc.Category{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GET /categories/{id}: Trae una categoría específica
func (s *Server) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	category, err := s.queries.GetCategory(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// PUT /categories/{id}: Modifica una categoría
func (s *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "El nombre de la categoría no puede estar vacío", http.StatusBadRequest)
		return
	}

	err = s.queries.UpdateCategory(r.Context(), sqlc.UpdateCategoryParams{
		ID:   int32(id),
		Name: input.Name,
	})

	if err != nil {
		http.Error(w, "Error al actualizar la categoría", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Categoría actualizada con éxito"}`))
}

// DELETE /categories/{id}: Elimina una categoría
func (s *Server) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	err = s.queries.DeleteCategory(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar la categoría", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
