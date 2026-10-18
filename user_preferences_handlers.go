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
func (s *Server) AddUserPreference(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     int32 `json:"user_id"`
		CategoryID int32 `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar la carga útil (JSON)", http.StatusBadRequest)
		return
	}

	err := s.queries.AddUserPreference(r.Context(), sqlc.AddUserPreferenceParams{
		UserID:     input.UserID,
		CategoryID: input.CategoryID,
	})
	if err != nil {
		http.Error(w, "Error de persistencia en la base de datos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Preferencia guardada con éxito"}`))
}

// GET /user-preferences/{user_id}
// Propósito: Obtiene y retorna todas las categorías asociadas a un identificador de usuario específico.
func (s *Server) ListUserPreferences(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		http.Error(w, "Identificador de usuario con formato inválido", http.StatusBadRequest)
		return
	}

	preferences, err := s.queries.ListUserPreferences(r.Context(), int32(userId))
	if err != nil {
		http.Error(w, "Error al consultar las preferencias", http.StatusInternalServerError)
		return
	}

	if preferences == nil {
		preferences = []sqlc.Category{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preferences)
}

// PUT /user-preferences/{user_id}
// Propósito: Reemplaza la totalidad de las preferencias de un usuario por un nuevo conjunto de categorías.
func (s *Server) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) RemoveUserPreference(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	catIdStr := r.PathValue("category_id")

	userId, err1 := strconv.ParseInt(userIdStr, 10, 32)
	catId, err2 := strconv.ParseInt(catIdStr, 10, 32)

	if err1 != nil || err2 != nil {
		http.Error(w, "Parámetros de URI inválidos", http.StatusBadRequest)
		return
	}

	err := s.queries.RemoveUserPreference(r.Context(), sqlc.RemoveUserPreferenceParams{
		UserID:     int32(userId),
		CategoryID: int32(catId),
	})
	if err != nil {
		http.Error(w, "Error de eliminación en la base de datos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
