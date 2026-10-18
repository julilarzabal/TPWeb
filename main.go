package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// Server actúa como contenedor de las dependencias necesarias para los handlers
// Decidimos utilizarlo para permitir que luego sea escalable y podamos agregar más dependencias si es necesario.
type Server struct {
	queries *sqlc.Queries
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que sea la ruta exacta "/" (para el 404)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Mostrar el formulario HTML
	http.ServeFile(w, r, "index.html")
}

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Estructura para recibir el cuerpo JSON de la petición
	/*
		Es una estructura anónima que se usa como un DTO (Data Transfer Object).
		Su única función es definir qué campos JSON esperas que el cliente
		(por ejemplo, Postman, un frontend o curl) te envíe en el cuerpo (body) de la petición HTTP.
	*/
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AgeRange    string `json:"age_range"`
		UserID      int32  `json:"user_id"`
	}

	// Carga de datos en req
	/*
		r.Body: Contiene el flujo de bytes (stream) con el JSON que envió el usuario.
		json.NewDecoder(...).Decode(&req): Lee ese JSON, analiza sus claves y rellena los campos de la variable req.
	*/
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Validación (ejemplo: validar que el título no esté vacío)
	if req.Title == "" {
		http.Error(w, "El título es obligatorio", http.StatusBadRequest)
		return
	}

	// 3. Ejecutar la consulta generada por sqlc
	/*
		s.queries.CreateEvent: Llama a la función generada por sqlc para crear un nuevo evento.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		db.CreateEventParams{...}: Estructura con los parámetros para la consulta.
	*/
	event, err := s.queries.CreateEvent(r.Context(), db.CreateEventParams{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		UserID:      req.UserID,
	})

	/*
		http.Error(): Es una función auxiliar de Go para devolver respuestas de error en texto plano. Recibe 3 parámetros:
		http.Error(w, "Mensaje de error", CÓDIGO_HTTP)
		(Ideal para responder errores)
	*/
	if err != nil {
		http.Error(w, "Error al crear el evento", http.StatusInternalServerError)
		return
	}

	// 4. Devolver la respuesta como JSON con estado 201 Created
	// w.WriteHeader(...): (Ideal para respuestas Exitosas)
	// a. Decirle al navegador/cliente que el contenido devuelto será JSON
	w.Header().Set("Content-Type", "application/json")

	// b. A AQUÍ SE INDICA EL CÓDIGO DE ESTADO (201 Created)
	w.WriteHeader(http.StatusCreated) // O directamente: w.WriteHeader(201)

	// c. Escribir el cuerpo del JSON que devolverás al cliente
	json.NewEncoder(w).Encode(event)
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
		events = []db.Event{}
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
		UserID      int32  `json:"user_id"`
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

	// 4. Ejecutar la actualización en sqlc
	/*
		s.queries.UpdateEvent: Llama a la función generada por sqlc para actualizar un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		db.UpdateEventParams{...}: Estructura con los parámetros para la consulta de actualización.
	*/
	err = s.queries.UpdateEvent(r.Context(), db.UpdateEventParams{
		ID:          int32(id),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		UserID:      req.UserID,
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

func main() {
	// Establecer la conexión a la base de datos
	connStr := "user=u password=p dbname=postgres host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Verificar que la base de datos responda
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}

	// Crear la instancia del repositorio generado por sqlc
	queries := sqlc.New(db)

	server := &Server{
		queries: queries,
	}

	// Configurar rutas utilizando las funciones del Server
	mux := http.NewServeMux()

	mux.HandleFunc("/", server.homeHandler)
	mux.HandleFunc("/events", server.CreateEventHandler)
	mux.HandleFunc("/events", server.ListEventsHandler)
	mux.HandleFunc("/events/{id}", server.GetEventHandler)
	mux.HandleFunc("/events/{id}", server.UpdateEventHandler)
	mux.HandleFunc("/events/{id}", server.DeleteEventHandler)
	mux.Handle("/style.css", http.FileServer(http.Dir(".")))

	fmt.Println("Servidor iniciado en http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}
