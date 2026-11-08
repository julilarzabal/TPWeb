package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parsear los datos del formulario clásico
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Extraer los datos usando el atributo 'name' de tus inputs en el .templ
	title := r.FormValue("title")
	descriptionStr := r.FormValue("description")
	location := r.FormValue("location")
	ageRange := r.FormValue("ageRange")
	creatorID, err := strconv.ParseInt(r.FormValue("creatorID"), 10, 32)
	if err != nil || creatorID <= 0 {
		http.Error(w, "El ID del organizador debe ser un entero positivo", http.StatusBadRequest)
		return
	}
	fechaStr := r.FormValue("fecha")

	// 2. Validación de regla de negocio
	if title == "" || fechaStr == "" || location == "" {
		http.Error(w, "Los campos 'title', 'fecha' y 'location' son obligatorios", http.StatusBadRequest)
		return
	}

	// 3. Convertir la fecha al tipo requerido por sqlc.
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "La fecha debe tener formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var description *string
	if descriptionStr != "" {
		description = &descriptionStr
	}

	// 4. Guardar en la base de datos mediante el método create de sqlc
	_, err = s.queries.CreateEvent(r.Context(), sqlc.CreateEventParams{
		Title:       title,
		Description: description,
		Location:    location,
		AgeRange:    ageRange,
		CreatorID:   int32(creatorID),
		Fecha:       fecha,
	})
	if err != nil {
		http.Error(w, "Error al guardar el evento en la base de datos", http.StatusInternalServerError)
		return
	}

	// 5. Redirigir al usuario (Patrón PRG)
	http.Redirect(w, r, "/", http.StatusSeeOther) // Esto es completamente correcto según el TP
}

func (s *Server) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener los registros mediante sqlc
	/*
		s.queries.ListEvents: Llama a la función generada por sqlc para obtener la lista de eventos.
		r.Context(): Pasa el contexto de la solicitud HTTP.
	*/
	events, err := s.queries.ListEvents(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de eventos", http.StatusInternalServerError)
		return
	}

	// Si no hay datos, devolver un array vacío en lugar de null
	if events == nil {
		events = []sqlc.Event{}
	}

	// 2. Responder en formato JSON
	w.Header().Set("Content-Type", "application/json")
	// No realiza w.WriteHeader(...) porque no debo devolver un código de error, sino un 200 OK por defecto.
	json.NewEncoder(w).Encode(events)
}

func (s *Server) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer el ID de la URL (usando net/http de Go)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Convertir el ID a entero
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Buscar en la base de datos con sqlc
	/*
		s.queries.GetEvent: Llama a la función generada por sqlc para obtener un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		int32(id): ID del evento a buscar.
	*/
	event, err := s.queries.GetEvent(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			// Si no existe, devolver estado 404
			http.Error(w, "Evento no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		return
	}

	// 3. Devolver el evento hallado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Convertir el ID a entero
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Leer JSON
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AgeRange    string `json:"age_range"`
		CreatorID   int32  `json:"creator_id"`
		Fecha       string `json:"fecha"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Validar datos
	if req.Title == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		http.Error(w, "La fecha es obligatoria y debe tener formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// 4. Ejecutar la actualización en sqlc
	/*
		s.queries.UpdateEvent: Llama a la función generada por sqlc para actualizar un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		sqlc.UpdateEventParams{...}: Estructura con los parámetros para la consulta de actualización.
	*/
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	err = s.queries.UpdateEvent(r.Context(), sqlc.UpdateEventParams{
		ID:          int32(id),
		Title:       req.Title,
		Description: description,
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		CreatorID:   req.CreatorID,
		Fecha:       fecha,
	})
	if err != nil {
		http.Error(w, "Error al actualizar el evento", http.StatusInternalServerError)
		return
	}

	// No es necesario devolver el evento actualizado, solo un estado 200 OK
	// Por eso no se realiza json.NewEncoder(w).Encode(...)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener ID
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Convertir el ID a entero
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Eliminar usando sqlc
	/*
		s.queries.DeleteEvent: Llama a la función generada por sqlc para eliminar un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		int32(id): ID del evento a eliminar.
	*/
	err = s.queries.DeleteEvent(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el evento", http.StatusInternalServerError)
		return
	}

	// 3. Responder con estado 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
