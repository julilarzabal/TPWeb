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
func (s *Server) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Definimos una estructura temporal para guardar el JSON que envía el usuario
	var input struct {
		Name string `json:"name"`
	}

	// Leemos el cuerpo (Body) de la petición HTTP y lo guardamos en 'input'
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Cumplimos la consigna: validamos que el nombre no esté vacío
	if input.Name == "" {
		http.Error(w, "El nombre de la categoría no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Llamamos al "cocinero" (sqlc) para que ejecute el INSERT en PostgreSQL
	category, err := s.queries.CreateCategory(r.Context(), input.Name)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Devolvemos el status 201 (Created) y el objeto recién guardado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GET /categories: Trae todas las categorías
func (s *Server) ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Ejecuta el SELECT de todas las categorías
	categories, err := s.queries.ListCategories(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las categorías", http.StatusInternalServerError)
		return
	}

	// Si la tabla está vacía, evitamos devolver 'null' y devolvemos un array vacío '[]'
	if categories == nil {
		categories = []sqlc.Category{}
	}

	// Empaquetamos la lista en JSON y la enviamos al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GET /categories/{id}: Trae una categoría específica
func (s *Server) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extraemos el número que el usuario puso en la URL (ej: /categories/5)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32) // Convertimos ese texto '5' al número entero 5
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Buscamos el registro en la base de datos por su ID
	category, err := s.queries.GetCategory(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			// Cumplimos la consigna: Si no existe, devolvemos error 404
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
func (s *Server) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenemos y convertimos el ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Leemos los nuevos datos enviados en el JSON
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Validamos que no nos manden un nombre vacío
	if input.Name == "" {
		http.Error(w, "El nombre de la categoría no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Ejecutamos el UPDATE usando la estructura de parámetros que armó sqlc
	err = s.queries.UpdateCategory(r.Context(), sqlc.UpdateCategoryParams{
		ID:   int32(id),
		Name: input.Name,
	})

	if err != nil {
		http.Error(w, "Error al actualizar la categoría", http.StatusInternalServerError)
		return
	}

	// Respondemos que todo salió bien (Status 200)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Categoría actualizada con éxito"}`))
}

// DELETE /categories/{id}: Elimina una categoría
func (s *Server) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenemos y convertimos el ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Ejecutamos el DELETE en la base de datos
	err = s.queries.DeleteCategory(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar la categoría", http.StatusInternalServerError)
		return
	}

	// Cumplimos la consigna: Devolver estado 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
