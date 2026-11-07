package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// HANDLERS USER_PREFERENCES
// POST /user-preferences
// Propósito: Crea un nuevo registro en la tabla intermedia vinculando un usuario con una categoría.
func (s *Server) AddUserPreferenceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     int32 `json:"user_id"`
		CategoryID int32 `json:"category_id"`
	}

	// Decodifica el cuerpo de la petición HTTP hacia la estructura de datos.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar la carga útil (JSON)", http.StatusBadRequest)
		return
	}

	// Invoca la capa de acceso a datos generada por sqlc para persistir la relación.
	err := s.queries.AddUserPreference(r.Context(), sqlc.AddUserPreferenceParams{
		UserID:     input.UserID,
		CategoryID: input.CategoryID,
	})
	if err != nil {
		http.Error(w, "Error de persistencia en la base de datos", http.StatusInternalServerError)
		return
	}

	// Retorna el código de estado HTTP 201 indicando la creación exitosa del recurso.
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Preferencia guardada con éxito"}`))
}

// GET /user-preferences/{user_id}
// Propósito: Obtiene y retorna todas las categorías asociadas a un identificador de usuario específico.
func (s *Server) ListUserPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	// Extrae y convierte el parámetro de la ruta a un tipo de dato numérico.
	userIdStr := r.PathValue("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		http.Error(w, "Identificador de usuario con formato inválido", http.StatusBadRequest)
		return
	}

	// Ejecuta la consulta de selección mediante la instancia de sqlc.
	preferences, err := s.queries.ListUserPreferences(r.Context(), int32(userId))
	if err != nil {
		http.Error(w, "Error al consultar las preferencias", http.StatusInternalServerError)
		return
	}

	// Previene la serialización de un valor nulo devolviendo una estructura vacía si no hay registros.
	if preferences == nil {
		preferences = []sqlc.Category{}
	}
	// Establece las cabeceras correspondientes y serializa la respuesta.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preferences)
}

// PUT /user-preferences/{user_id}
// Propósito: Reemplaza la totalidad de las preferencias de un usuario por un nuevo conjunto de categorías.
func (s *Server) UpdateUserPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		http.Error(w, "Identificador de usuario con formato inválido", http.StatusBadRequest)
		return
	}

	var input struct {
		CategoryIDs []int32 `json:"category_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar la carga útil (JSON)", http.StatusBadRequest)
		return
	}

	// Ejecuta la transacción de actualización masiva definida en la capa de datos.
	err = s.queries.UpdateUserPreferences(r.Context(), sqlc.UpdateUserPreferencesParams{
		UserID:  int32(userId),
		Column2: input.CategoryIDs,
	})
	if err != nil {
		http.Error(w, "Error al actualizar la colección de preferencias", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Preferencias actualizadas con éxito"}`))
}

// DELETE /user-preferences/{user_id}/{category_id}
// Propósito: Elimina un registro específico en la tabla intermedia dados ambos identificadores.
func (s *Server) RemoveUserPreferenceHandler(w http.ResponseWriter, r *http.Request) {
	// Extrae secuencialmente ambos parámetros requeridos de la URI.
	userIdStr := r.PathValue("user_id")
	catIdStr := r.PathValue("category_id")

	userId, err1 := strconv.ParseInt(userIdStr, 10, 32)
	catId, err2 := strconv.ParseInt(catIdStr, 10, 32)

	if err1 != nil || err2 != nil {
		http.Error(w, "Parámetros de URI inválidos", http.StatusBadRequest)
		return
	}

	// Ejecuta la sentencia de eliminación mediante sqlc.
	err := s.queries.RemoveUserPreference(r.Context(), sqlc.RemoveUserPreferenceParams{
		UserID:     int32(userId),
		CategoryID: int32(catId),
	})
	if err != nil {
		http.Error(w, "Error de eliminación en la base de datos", http.StatusInternalServerError)
		return
	}

	// Retorna el código de estado HTTP 204 indicando éxito sin contenido en el cuerpo.
	w.WriteHeader(http.StatusNoContent)
}
