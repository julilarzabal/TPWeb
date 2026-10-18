package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// HANDLERS DE EVENT_CATEGORY

/*
Dado que event_category representa una relación entre un Evento y una Categoría,
los endpoints REST usualmente se diseñan como sub-recursos de /events/{id}/categories.
*/

// En este handler, obtenemos el event_id desde la URL y el category_id desde el cuerpo JSON de la petición.
func (s *Server) AddEventCategoryHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		CategoryID int32 `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.CategoryID <= 0 {
		http.Error(w, "El category_id es obligatorio y debe ser mayor a 0", http.StatusBadRequest)
		return
	}

	err = s.queries.AddEventCategory(r.Context(), db.AddEventCategoryParams{
		EventID:    int32(eventID),
		CategoryID: req.CategoryID,
	})
	if err != nil {
		http.Error(w, "Error al asociar la categoría al evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ListEventCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	categories, err := s.queries.ListEventCategories(r.Context(), int32(eventID))
	if err != nil {
		http.Error(w, "Error al consultar las categorías del evento", http.StatusInternalServerError)
		return
	}

	if categories == nil {
		categories = []db.Category{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Recibe ambos IDs por la URL ({id} para el evento y {category_id} para la categoría) y elimina únicamente ese vínculo usando RemoveEventCategory.
func (s *Server) RemoveEventCategoryHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	categoryIDStr := r.PathValue("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	err = s.queries.RemoveEventCategory(r.Context(), db.RemoveEventCategoryParams{
		EventID:    int32(eventID),
		CategoryID: int32(categoryID),
	})
	if err != nil {
		http.Error(w, "Error al remover la categoría del evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Utiliza la consulta DeleteEventCategories para desvincular todas las categorías de un evento a la vez
func (s *Server) DeleteEventCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	err = s.queries.DeleteEventCategories(r.Context(), int32(eventID))
	if err != nil {
		http.Error(w, "Error al eliminar las categorías del evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
