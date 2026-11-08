package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// CreateCategoryHandler procesa el POST del formulario
func (s *Server) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parsear los datos del formulario clásico
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Extraer los datos usando el atributo 'name' de tus inputs en el .templ
	name := r.FormValue("name")

	// 2. Validación de regla de negocio
	if name == "" {
		http.Error(w, "El campo 'name' es obligatorio", http.StatusBadRequest)
		return
	}

	// 4. Guardar en la base de datos mediante el método create de sqlc
	_, err := s.queries.CreateCategory(r.Context(), name)
	if err != nil {
		http.Error(w, "Error al guardar la categoria en la base de datos", http.StatusInternalServerError)
		return
	}

	// 5. Redirigir al usuario (Patrón PRG)
	http.Redirect(w, r, "/", http.StatusSeeOther) // Esto es completamente correcto según el TP
}

/*
	HANDLERS TP3
*/

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
