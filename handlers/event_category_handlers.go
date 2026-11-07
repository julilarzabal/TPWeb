package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

/*
Dado que event_category representa una relación entre un Evento y una Categoría,
los endpoints REST usualmente se diseñan como sub-recursos de /events/{id}/categories.
*/

// En este handler, obtenemos el event_id desde la URL y el category_id desde el cuerpo JSON de la petición.
func (s *Server) AddEventCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el ID del evento desde la URL
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr) // Convertir el ID a entero
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	// 2. Leer el cuerpo JSON que contiene el category_id
	var req struct {
		CategoryID int32 `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Validar que el category_id sea válido
	if req.CategoryID <= 0 {
		http.Error(w, "El category_id es obligatorio y debe ser mayor a 0", http.StatusBadRequest)
		return
	}

	// 4. Ejecutar la consulta de sqlc
	/*
		s.queries.AddEventCategory: Llama a la función generada por sqlc para asociar una categoría a un evento.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		db.AddEventCategoryParams{...}: Parámetros para la consulta.
	*/
	err = s.queries.AddEventCategory(r.Context(), db.AddEventCategoryParams{
		EventID:    int32(eventID),
		CategoryID: req.CategoryID,
	})
	if err != nil {
		http.Error(w, "Error al asociar la categoría al evento", http.StatusInternalServerError)
		return
	}

	// 5. Responder con estado 201 Created
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ListEventCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el ID del evento desde la URL
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr) // Convertir el ID a entero
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	// 2. Ejecutar la consulta de sqlc
	/*
		s.queries.ListEventCategories: Llama a la función generada por sqlc para obtener las categorías de un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		int32(eventID): ID del evento del cual se quieren obtener las categorías.
	*/
	categories, err := s.queries.ListEventCategories(r.Context(), int32(eventID))
	if err != nil {
		http.Error(w, "Error al consultar las categorías del evento", http.StatusInternalServerError)
		return
	}

	// Evitar devolver null en JSON si no hay categorías asociadas
	if categories == nil {
		categories = []db.Category{}
	}

	// 3. Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Recibe ambos IDs por la URL ({id} para el evento y {category_id} para la categoría) y elimina únicamente ese vínculo usando RemoveEventCategory.
func (s *Server) RemoveEventCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener los IDs de la URL
	// ID del evento
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr) // Convertir el ID a entero
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	// ID de la categoría
	categoryIDStr := r.PathValue("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// 2. Ejecutar el borrado en la base de datos con sqlc
	/*
		s.queries.RemoveEventCategory: Llama a la función generada por sqlc para remover una categoría de un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		db.RemoveEventCategoryParams{...}: Parámetros para la consulta.
	*/
	err = s.queries.RemoveEventCategory(r.Context(), db.RemoveEventCategoryParams{
		EventID:    int32(eventID),
		CategoryID: int32(categoryID),
	})
	if err != nil {
		http.Error(w, "Error al remover la categoría del evento", http.StatusInternalServerError)
		return
	}

	// 3. Responder con 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// Utiliza la consulta DeleteEventCategories para desvincular todas las categorías de un evento a la vez
func (s *Server) DeleteEventCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener ID del evento
	eventIDStr := r.PathValue("id")
	eventID, err := strconv.Atoi(eventIDStr) // Convertir el ID a entero
	if err != nil || eventID <= 0 {
		http.Error(w, "ID de evento inválido", http.StatusBadRequest)
		return
	}

	// 2. Ejecutar borrado masivo en sqlc
	/*
		s.queries.DeleteEventCategories: Llama a la función generada por sqlc para eliminar todas las categorías de un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		int32(eventID): ID del evento del cual se quieren eliminar las categorías.
	*/
	err = s.queries.DeleteEventCategories(r.Context(), int32(eventID))
	if err != nil {
		http.Error(w, "Error al eliminar las categorías del evento", http.StatusInternalServerError)
		return
	}

	// 3. Responder con 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
