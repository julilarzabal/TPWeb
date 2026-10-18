package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	db "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// HANDLERS DE EVENT

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AgeRange    string `json:"age_range"`
		UserID      int32  `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "El título es obligatorio", http.StatusBadRequest)
		return
	}

	event, err := s.queries.CreateEvent(r.Context(), db.CreateEventParams{
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		UserID:      req.UserID,
	})
	if err != nil {
		http.Error(w, "Error al crear el evento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func (s *Server) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := s.queries.ListEvents(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de eventos", http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []db.Event{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	event, err := s.queries.GetEvent(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Evento no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AgeRange    string `json:"age_range"`
		UserID      int32  `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	err = s.queries.UpdateEvent(r.Context(), db.UpdateEventParams{
		ID:          int32(id),
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		UserID:      req.UserID,
	})
	if err != nil {
		http.Error(w, "Error al actualizar el evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = s.queries.DeleteEvent(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
